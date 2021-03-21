package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go-speedtest-bot/module/controller"
	"sync"
)

type B = tgbotapi.BotAPI
type M = tgbotapi.Message

type CMDFunc func(*M)
type CMD = map[string]CMDFunc

var Commands = CMD{
	"start":      cmdStart,
	"ping":       cmdPing,
	"status":     cmdStatus,
	"read_sub":   cmdReadSub,
	"result":     cmdResult,
	"run_url":    cmdStartTestWithURL,
	"list_subs":  cmdListSubs,
	"schedule":   cmdSchedule,
	"listrunner": cmdListRunner,
}

// DefaultConfig contains all the default speed test setting.
type DefaultConfig struct {
	Remarks  string
	Url      string
	Method   string
	Mode     string
	Interval int
	Chat     int64
	Include  []string
	Exclude  []string
}

type SpeedTestArguments struct {
	Subs    string   `json:"subs"`
	Method  string   `json:"method"`
	Mode    string   `json:"mode"`
	Exclude []string `json:"exclude"`
	Include []string `json:"include"`
}

type RunnerComm struct {
	m map[string]*controller.Comm
	s sync.Mutex
}

func (rc *RunnerComm) Exist(n string) bool {
	rc.s.Lock()
	defer rc.s.Unlock()
	_, ok := rc.m[n]
	return ok
}

func (rc *RunnerComm) Register(n string, c *controller.Comm) {
	rc.s.Lock()
	defer rc.s.Unlock()
	rc.m[n] = c
}

func (rc *RunnerComm) UnRegister(n string) {
	rc.s.Lock()
	defer rc.s.Unlock()
	if _, ok := rc.m[n]; ok {
		delete(rc.m, n)
	}
}

func (rc *RunnerComm) C(n string) *controller.Comm {
	rc.s.Lock()
	defer rc.s.Unlock()
	if c, ok := rc.m[n]; ok {
		return c
	}
	return nil
}
