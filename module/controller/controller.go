package controller

import (
	"errors"
	"fmt"
	"go-speedtest-bot/module/config"
	"go-speedtest-bot/module/heartbeat"
	"go-speedtest-bot/module/runner"
	"go-speedtest-bot/module/speedtest"
	"time"
)

var (
	runnerNotFoundErr = errors.New("runner not found")
	permDenied        = errors.New("permission denied")
)

// Ping test the connection of the given runner
func Ping(requester int, r *runner.Runner) error {
	if !r.HasAccess(requester) {
		return permDenied
	}
	connected := speedtest.Ping(*r)
	if connected {
		return nil
	}
	return errors.New("unable to connect to the " + r.Name + " host")
}

// PingFromName using given name to test runner connection.
// Return error if given runner name not found.
func PingFromName(requester int, name string) error {
	r := config.GetRunner(name)
	if r == nil {
		return runnerNotFoundErr
	}
	return Ping(requester, r)
}

// GetStatus return runner backend status. If given runner
// name not found, or backend has error, return null string
// and error.
func GetStatus(requester int, name string) (string, error) {
	r := config.GetRunner(name)
	if r == nil {
		return "", runnerNotFoundErr
	}

	if !r.HasAccess(requester) {
		return "", permDenied
	}

	result, err := speedtest.GetStatus(*r)
	if err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", errors.New(result.Error)
	}

	return fmt.Sprintf("%s status: %s", name, result.State), nil
}

// ListRunner return a list of runner that the requester can
// access.
func ListRunner(requester int) string {
	runners := config.GetAllRunner()
	var text string = "available runner:\n"
	for _, r := range runners {
		if r.HasAccess(requester) {
			text += r.Name + "\n"
		}
	}
	return text
}

// ReadSubscriptions parse a subscriptions url to readable text
func ReadSubscriptions(requester int, sub, name string) (string, error) {
	r := config.GetRunner(name)
	if r == nil {
		return "", runnerNotFoundErr
	}
	if !r.HasAccess(requester) {
		return "", permDenied
	}

	resp, err := speedtest.ReadSubscriptions(*r, sub)
	if err != nil {
		return "", err
	}
	if len(resp) == 0 {
		return "", errors.New("unknown subscription")
	}
	var text = resp[0].Type
	for _, subResp := range resp {
		if subResp.Config != nil {
			text += subResp.Config.Remarks + "\n"
		}
	}
	return text, nil
}

// Result return the test result of given runner
func Result(requester int, name string) (string, error) {
	r := config.GetRunner(name)
	if r == nil {
		return "", runnerNotFoundErr
	}

	if !r.HasAccess(requester) {
		return "", permDenied
	}

	result, err := speedtest.GetResult(*r)
	if err != nil {
		return "", err
	}

	return formatResult(result), nil
}

func formatResult(r *speedtest.Result) string {
	if len(r.Result) == 0 {
		return "no result"
	}

	var text string
	text = "Recent result:\n"
	text += "Runner Status: " + r.Status + "\n"

	if len(r.Current.Remarks) != 0 {
		text += "Nodes being tested: " + r.Current.Remarks + "\n"
	}

	text += "\n"
	for _, res := range r.Result {
		text += fmt.Sprintf(
			"%s: | loss: %.2f%% | local ping: %.2f ms | google ping: %.2f ms\n",
			res.Remarks, res.Loss*100, res.Ping*1000, res.GPing*1000,
		)
	}
	return text
}

// Run run a test with given configuration
func Run(requester int, name, sub, method, mode string, include, exclude []string, errCh chan error) error {
	r := config.GetRunner(name)
	if r == nil {
		return runnerNotFoundErr
	}
	if !r.HasAccess(requester) {
		return permDenied
	}
	if r.IsWorking() {
		return errors.New("backend is running other jobs")
	}

	nodes, err := speedtest.ReadSubscriptions(*r, sub)
	if err != nil {
		return fmt.Errorf("%s reading %s: %v", name, sub, err)
	}

	if sub == "" {
		return errors.New("no subscriptions specific")
	}
	if method == "" {
		return errors.New("no method specific")
	}
	if mode == "" {
		return errors.New("no mode specific")
	}
	if len(include) != 0 {
		nodes = speedtest.IncludeRemarks(nodes, include)
	}
	if len(exclude) != 0 {
		nodes = speedtest.ExcludeRemarks(nodes, exclude)
	}

	cfg := speedtest.NewStartConfigs(method, mode, nodes)

	if errCh == nil {
		return errors.New("can't pass error in nil channel")
	}

	go func() {
		msg, err := speedtest.StartTest(*r, cfg)
		if err != nil {
			errCh <- err
			return
		}

		if msg == "running" {
			errCh <- errors.New("backend running other jobs")
			return
		}

		errCh <- nil
	}()

	return nil
}

// Schedule is like cron jobs, running test with user predefine configuration
func Schedule(requester int, cfg *config.Default, c *Comm) error {
	r := config.GetRunner(cfg.DefaultRunner)
	if r == nil {
		return runnerNotFoundErr
	}
	if !r.HasAccess(requester) {
		return permDenied
	}
	if cfg == nil {
		return errors.New("configuration are empty")
	}
	if c.ErrCh == nil || c.LogCh == nil {
		return errors.New("no communicate channel opened")
	}
	go schedule(r, cfg, c)
	return nil
}

func schedule(r *runner.Runner, cfg *config.Default, c *Comm) {
	heartBeat := heartbeat.NewHeartBeat(cfg.Interval)
	resultCh := make(chan string)
	for {
		select {
		case <-heartBeat.C():
			tc, err := newTestConfig(r, cfg)
			if err != nil {
				c.ErrCh <- err
				continue
			}

			heartBeat.Stop()

			go startScheduleTest(r, tc, c, heartBeat, resultCh)

			c.LogCh <- r.Name + " start new jobs"

		case state := <-resultCh:
			c.LogCh <- "runner " + r.Name + " finish one test"
			if state == "running" {
				c.ErrCh <- errors.New("another jobs running, please start schedule jobs later")
				heartBeat.Reset()
				return
			}

			re, err := speedtest.GetResult(*r)
			if err != nil {
				c.ErrCh <- fmt.Errorf("schedule get result: %v", err)
				heartBeat.Reset()
				return
			}

			alert := AlertHandler(re.Result)
			c.Alert <- &alert

			heartBeat.Reset()

		case <-c.Sig:
			c.LogCh <- r.Name + " exit schedule jobs"
			return
		}
	}
}

func newTestConfig(r *runner.Runner, cfg *config.Default) (*speedtest.StartConfigs, error) {
	var retry int
	nodes, err := speedtest.ReadSubscriptions(*r, cfg.Link)
	for err != nil {
		if retry > 5 {
			return nil, fmt.Errorf("retry after 5 times: %v", err)
		}
		time.Sleep(5 * time.Second)
		retry++
		nodes, err = speedtest.ReadSubscriptions(*r, cfg.Link)
	}
	return speedtest.NewStartConfigs("ST_ASYNC", "TCP_PING", nodes), nil
}

func startScheduleTest(r *runner.Runner, tc *speedtest.StartConfigs, c *Comm, beat *heartbeat.HeartBeat, resultCh chan string) {
	resp, err := speedtest.StartTest(*r, tc)
	if err != nil {
		c.ErrCh <- fmt.Errorf("%s run schedule test: %v", r.Name, err)
		c.LogCh <- r.Name + " schedule job got error, recovered"
		beat.Reset()
		return
	}
	resultCh <- resp
}
