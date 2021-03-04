package config

import "sync/atomic"

// Configuration contains bot, runner configuration
// at run time; Also store user define default
// subscription config
type Configuration struct {
	Global        Global     `json:"global"`
	Runner        []*Runner  `json:"runner"`
	DefaultConfig []*Default `json:"default_config"`
}

type Global struct {
	Token string `json:"token"`
	Admin int    `json:"admin"`
}

type Runner struct {
	Status int32  // Status store runner status at local
	Name   string `json:"name"`
	Host   Host   `json:"host"`
	Admins []int  `json:"admins"`
}

const (
	Pending = iota
	Working
)

// GetRunnerStatus return current status
// 0 == Pending
// 1 == Working
func (r *Runner) GetRunnerStatus() int32 {
	return atomic.LoadInt32(&r.Status)
}

// HangUp changed runner status to pending
func (r *Runner) HangUp() {
	atomic.CompareAndSwapInt32(&r.Status, Working, Pending)
}

// Activate changed runner status to working
func (r *Runner) Activate() {
	atomic.CompareAndSwapInt32(&r.Status, Pending, Working)
}

type Host struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

func (h *Host) GetURL() string {
	return h.Address
}

type Default struct {
	Name          string `json:"name"`
	Link          string `json:"link"`
	Chat          int64  `json:"chat"`
	Admins        []int  `json:"admins"`
	DefaultRunner string `json:"default_runner"`
}
