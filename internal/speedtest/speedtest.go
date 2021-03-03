package speedtest

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-speedtest-bot/internal/web"
	"os"
	"path"
	"strings"
)

// Ping will test if the given host is accessible or not
func (r *Runner) Ping() bool {
	resp, err := web.Get(path.Join(r.Host.GetURL(), "getversion"))
	if err != nil {
		return false
	}
	var v Version
	err = json.Unmarshal(resp, &v)
	if err != nil {
		return false
	}

	return v.Main != "" && v.WebAPI != ""
}

// GetStatus is used for fetching backend status
func (r *Runner) GetStatus() (*Status, error) {
	resp, err := web.Get(path.Join(r.Host.GetURL(), "status"))
	if err != nil {
		return nil, fmt.Errorf("get status: %v", err)
	}
	var st *Status
	err = json.Unmarshal(resp, &st)
	if err != nil {
		return nil, fmt.Errorf("decode %q: %v", resp, err)
	}
	return st, nil
}

// ReadSubscriptions return list of node information with the given subscription url.
func (r *Runner) ReadSubscriptions(sub string) ([]*SubscriptionResp, error) {
	data := map[string]string{"url": sub}
	jsondata, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("%s not valid", sub)
	}
	resp, err := web.JSONPost(path.Join(r.Host.GetURL(), "readsubscriptions"), jsondata)
	if err != nil {
		return nil, fmt.Errorf("post sub: %v", err)
	}
	var cfg []*SubscriptionResp
	err = json.Unmarshal(resp, &cfg)
	if err != nil {
		return nil, fmt.Errorf("%q not valid", resp)
	}
	return cfg, nil
}

// GetResult return the newest speed test result.
func (r *Runner) GetResult() (*Result, error) {
	resp, err := web.Get(path.Join(r.Host.GetURL(), "getresults"))
	if err != nil {
		return nil, fmt.Errorf("get result: %v", err)
	}
	var result Result
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return nil, fmt.Errorf("%q not valid", resp)
	}
	return &result, nil
}

// StartTest post a speed test request with the given node config.
// Because of the blocking backend, This function is designed to be called as goroutine.
// If error happen function will pass error message into status chan.
// If state is not null it will pass status: ["running" / "done"].
// If you are calling this method without authorized of given unexpected config it will pass error message.
func (r *Runner) StartTest(startCFG *StartConfigs, statusChan chan string) {
	d, err := json.Marshal(startCFG)
	if err != nil {
		e := sendStatus(statusChan, "invalid start config")
		if e != nil {
			fmt.Println(e)
			os.Exit(0)
		}
		return
	}
	resp, err := web.JSONPostWithTimeout(path.Join(r.Host.GetURL(), "start"), d, 0)
	if err != nil {
		sendStatus(statusChan, fmt.Sprintf("speed test failed, response: %q", resp))
		return
	}
	var state Status
	err = json.Unmarshal(resp, &state)
	if err != nil {
		sendStatus(statusChan, fmt.Sprintf(""))
		return
	}
	if state.State != "" {
		sendStatus(statusChan, state.State)
	} else {
		sendStatus(statusChan, state.Error)
	}
}

func sendStatus(status chan string, content string) error {
	var err error
	if status == nil {
		return errors.New("nil status channel")
	}
	defer func() {
		msg := recover()
		err = fmt.Errorf("send status: %v", msg)
	}()
	status <- content
	return err
}

// IncludeRemarks will select all the configs with the given remarks.
func IncludeRemarks(configs []*SubscriptionResp, incRems []string) []*SubscriptionResp {
	var newcfg []*SubscriptionResp
	for _, e := range incRems {
		for _, c := range configs {
			if strings.Contains(c.Config.Remarks, e) {
				newcfg = append(newcfg, c)
				break
			}
		}
	}
	return newcfg
}

// ExcludeRemarks will select all the configs excluded the given remarks.
func ExcludeRemarks(configs []*SubscriptionResp, excludeRemarks []string) []*SubscriptionResp {
	var newcfg []*SubscriptionResp
	remarksLength := len(excludeRemarks)
	for _, c := range configs {
		for i, mark := range excludeRemarks {
			if mark == c.Config.Remarks {
				break
			}
			if i == remarksLength-1 {
				newcfg = append(newcfg, c)
			}
		}
	}
	return newcfg
}
