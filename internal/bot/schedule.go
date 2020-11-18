package bot

import (
	"go-speedtest-bot/internal/speedtest"
	"log"
	"os"
	"time"
)

var pause bool

func start(b *B) error {
	if pause {
		pause = false
	}
	log.Println("[Schedule]New loop started")
	request(b)
	return nil
}

func fetchResult() string {
	result, err := speedtest.GetResult(speedtest.GetHost())
	if err != nil {
		return err.Error()
	}
	return formatResult(result)
}

func request(b *B) {
	nodes, err := speedtest.ReadSubscriptions(speedtest.GetHost(), Def.Url)
	if err != nil {
		pause = true
		log.Println("[ReadSubscriptionsError]", err)
		SendT(b, Def.Chat, err.Error())
		return
	}

	cfg := speedtest.NewStartConfigs(Def.Method, Def.Mode, nodes)
	host := speedtest.GetHost()
	stChan := make(chan string)
	go speedtest.StartTest(host, cfg, stChan)
	select {
	case s := <-stChan:
		if s == "done" {
			SendT(b, Def.Chat, fetchResult())
			wait(Def.Interval)
			go speedtest.StartTest(host, cfg, stChan)
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
