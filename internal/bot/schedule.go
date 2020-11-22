package bot

import (
	"go-speedtest-bot/internal/speedtest"
	"log"
	"os"
	"time"
)

var pause, started, alert bool

func start(b *B) {
	if pause {
		pause = false
	}
	log.Println("[Schedule]New loop started")
	started = true
	request(b)
	log.Println("[Schedule]loop stopped")
}

func fetchResult() []speedtest.ResultInfo {
	result, err := speedtest.GetResult(speedtest.GetHost())
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return result.Result
}

func request(b *B) {
	nodes, err := speedtest.ReadSubscriptions(speedtest.GetHost(), Def.Url)
	if err != nil {
		pause = true
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
	for {
		go speedtest.StartTest(host, cfg, stChan)
		select {
		case s := <-stChan:
			if s == "done" {
				AlertHandler(fetchResult(), b)
				wait(Def.Interval)
			} else {
				pause = true
				log.Println("[SpeedTestError]" + s)
				SendT(b, Def.Chat, s)
			}
			if pause {
				return
			}
		}
	}
}

func wait(s int) {
	if s < 1 {
		log.Println("[Sleep]Unexpected interval")
		os.Exit(-1)
	}
	time.Sleep(time.Second * time.Duration(s))
}

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
