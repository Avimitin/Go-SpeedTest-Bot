package bot

import (
	"go-speedtest-bot/internal/speedtest"
	"log"
	"sync/atomic"
	"time"
)

const (
	STOPPED int32 = iota
	RUNNING
)

type Job struct {
	status int32
	stop   chan int32
}

func NewJob() *Job {
	return &Job{
		status: STOPPED,
		stop:   make(chan int32, 1),
	}
}

func (j *Job) running() {
	ok := atomic.CompareAndSwapInt32(&j.status, STOPPED, RUNNING)
	if ok {
		log.Println("[INFO]Schedule job has started.")
	}
}

func (j *Job) stopped() {
	ok := atomic.CompareAndSwapInt32(&j.status, RUNNING, STOPPED)
	if ok {
		log.Println("[INFO]Schedule job has stopped.")
	}
}

func (j *Job) start(b *B) {
	log.Println("[Schedule]New goroutine started")
	atomic.CompareAndSwapInt32(&j.status, STOPPED, RUNNING)
	go j.Run(b)
}

func (j *Job) Run(b *B) {
	getTestCFG := func() *speedtest.StartConfigs {
		nodes, err := speedtest.ReadSubscriptions(speedtest.GetHost(), Def.Url)
		if err != nil {
			j.stopped()
			log.Println("[ReadSubscriptionsError]", err)
			SendT(b, Def.Chat, err.Error())
			return nil
		}
		if len(Def.Include) != 0 {
			nodes = speedtest.IncludeRemarks(nodes, Def.Include)
		}
		if len(Def.Exclude) != 0 {
			nodes = speedtest.ExcludeRemarks(nodes, Def.Exclude)
		}
		return speedtest.NewStartConfigs(Def.Method, Def.Mode, nodes)
	}
	host := speedtest.GetHost()
	stChan := make(chan string)
	period := time.Duration(Def.Interval) * time.Second
	t := time.NewTicker(period)
	for {
		select {
		case <-t.C:
			log.Println("[Schedule]New test started")
			cfg := getTestCFG()
			if cfg == nil {
				j.stopped()
				log.Println("[Schedule]Schedule jobs exit: get nil config.")
				return
			}
			go speedtest.StartTest(host, cfg, stChan)
			t.Stop()
		case s := <-stChan:
			if s == "done" {
				AlertHandler(fetchResult(), b)
				t.Reset(period)
			} else {
				log.Println("[SpeedTestError]" + s)
				SendT(b, Def.Chat, s+" Please restart schedule jobs.")
				j.stopped()
				return
			}
		case state := <-j.stop:
			switch state {
			case -1:
				log.Println("[Schedule]Jobs exit unexpectedly")
			case 0:
				log.Println("[Schedule]loop stopped")
			}
			return
		}
	}
}

func (j *Job) Stop(state int32) {
	j.stopped()
	j.stop <- state
}

func SetInterval(i int) {
	Def.Interval = i
}

func fetchResult() []speedtest.ResultInfo {
	result, err := speedtest.GetResult(speedtest.GetHost())
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return result.Result
}
