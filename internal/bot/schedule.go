package bot

import (
	"fmt"
	"go-speedtest-bot/internal/config"
	"go-speedtest-bot/internal/runner"
	"go-speedtest-bot/internal/speedtest"
	"log"
	"time"
)

func StartScheduleJobs(runner *runner.Runner, subsConfig *config.Default) {
	runner.Activate()
	period := (time.Duration(subsConfig.Interval) * time.Second) / 2
	heartBeat := time.NewTicker(period)
	runnerCh := runner.NewChan()
	statusCh := make(chan string)
	for {
		select {
		case <-heartBeat.C:
			var err error
			var startConfigs *speedtest.StartConfigs
			startConfigs, err = newTestConfig(runner, subsConfig)
			if err != nil {
				log.Println(err)
				SendTF(subsConfig.Chat, "schedule job: %v", err)
				return
			}
			heartBeat.Stop()
			go speedtest.StartTest(*runner, startConfigs, statusCh)
			log.Printf("runner %s start new test", runner.Name)
		case state := <-statusCh:
			if state == "done" {
				log.Printf("runner %s finish one test", runner.Name)
				handleResult(runner)
			} else {
				log.Printf("speedtest error: %s", state)
				SendT(subsConfig.Chat, "speedtest error: "+state+" \nSchedule job exit.")
				runner.HangUp()
				return
			}
			heartBeat.Reset(period)
		case <-runnerCh:
			runner.HangUp()
			return
		}
	}
}

func newTestConfig(runner *runner.Runner, subsConfig *config.Default) (*speedtest.StartConfigs, error) {
	var retry int
	nodes, err := speedtest.ReadSubscriptions(*runner, subsConfig.Link)
	for err != nil {
		if retry > 5 {
			return nil, fmt.Errorf("retry after 5 times: %v", err)
		}
		nodes, err = speedtest.ReadSubscriptions(*runner, subsConfig.Link)
	}
	return speedtest.NewStartConfigs("ST_ASYNC", "TCP_PING", nodes), nil
}

func handleResult(r *runner.Runner) {

}
