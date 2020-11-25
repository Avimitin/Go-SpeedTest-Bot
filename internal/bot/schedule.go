package bot

import (
	"go-speedtest-bot/internal/speedtest"
	"log"
	"sync"
	"time"
)

type Job struct {
	started bool
	stop    chan bool
	mtx     sync.Mutex
}

var task *Job = &Job{
	started: false,
}

func start(b *B) {
	log.Println("[Schedule]New goroutine started")
	task.started = true
	go func() {
		task.Run(b)
	}()
}

func fetchResult() []speedtest.ResultInfo {
	result, err := speedtest.GetResult(speedtest.GetHost())
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return result.Result
}

func (j *Job) Run(b *B) {
	nodes, err := speedtest.ReadSubscriptions(speedtest.GetHost(), Def.Url)
	if err != nil {
		task.started = false
		log.Println("[ReadSubscriptionsError]", err)
		SendT(b, Def.Chat, err.Error())
		return
	}
	if len(Def.Include) != 0 {
		nodes = speedtest.IncludeRemarks(nodes, Def.Include)
	}
	if len(Def.Exclude) != 0 {
		nodes = speedtest.ExcludeRemarks(nodes, Def.Exclude)
	}
	cfg := speedtest.NewStartConfigs(Def.Method, Def.Mode, nodes)
	host := speedtest.GetHost()
	stChan := make(chan string)
	period := time.Duration(Def.Interval) * time.Second
	t := time.NewTicker(period)
	for {
		select {
		case <-t.C:
			log.Println("[Schedule]New test started")
			go speedtest.StartTest(host, cfg, stChan)
			t.Stop()
		case s := <-stChan:
			if s == "done" {
				AlertHandler(fetchResult(), b)
				t.Reset(period)
			} else {
				j.Stop()
				log.Println("[SpeedTestError]" + s)
				SendT(b, Def.Chat, s)
				return
			}
		case <-j.stop:
			log.Println("[Schedule]loop stopped")
			return
		}
	}
}

func (j *Job) Stop() {
	j.mtx.Lock()
	defer j.mtx.Unlock()
	j.started = false
	j.stop <- true
}

/*
func wait(s int) {
	if s < 1 {
		log.Println("[Sleep]Unexpected interval")
		os.Exit(-1)
	}
	time.Sleep(time.Second * time.Duration(s))
}
*/

func SetInterval(i int) {
	Def.Interval = i
}

//func longPoll(status chan string) {
//	state, err := speedtest.GetStatus(speedtest.GetHost())
//	if err != nil {
//		status <- err.Error()
//		return
//	}
//	status <- state.State
//}

//func startLoop(b *B) {
//	state := make(chan string)
//	go longPoll(state)
//
//	select {
//	case s := <-state:
//		if s == "done" {
//			SendT(b, Def.Chat, fetchResult())
//			wait(Def.Interval)
//			err := request()
//			if err != nil {
//				SendT(b, Def.Chat, "Task cancelled due to unexpected error"+err.Error())
//				return
//			}
//		} else if s == "running" {
//			wait(1)
//			go longPoll(state)
//		} else {
//			SendT(b, Def.Chat, "Timed task cancelled due to error "+s)
//			return
//		}
//
//		if pause {
//			SendT(b, Def.Chat, "Task pause.")
//			return
//		}
//	}
//}
