package speedtest

import (
	"sync"
	"sync/atomic"
)

const (
	Pending = iota
	Working
)

// Host determine Runner's remote address
type Host interface {
	GetURL() string
}

// Runner store backend address and status
type Runner struct {
	mu     sync.RWMutex
	Status int32
	Name   string
	Host   Host
	Admin  []int
}

// NewRunner return a new pointer to Runner
func NewRunner(name string, host Host, admin []int) *Runner {
	return &Runner{
		Status: Pending,
		Host:   host,
		Name:   name,
		Admin:  admin,
	}
}

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

type Version struct {
	Main   string `json:"main"`
	WebAPI string `json:"webapi"`
}

type Status struct {
	State string `json:"status"`
	Error string `json:"error"`
}

type SubscriptionResp struct {
	Type   string        `json:"type"`
	Config *ShadowConfig `json:"config"`
	Error  string        `json:"-"`
}

type ShadowConfig struct {
	Server        string `json:"server"`
	ServerPort    int    `json:"server_port"`
	Method        string `json:"method"`
	Protocol      string `json:"protocol"`
	Obfs          string `json:"obfs"`
	Plugin        string `json:"plugin"`
	Password      string `json:"password"`
	ProtocolParam string `json:"protocol_param"`
	ObfsParam     string `json:"obfsparam"`
	PluginOpts    string `json:"plugin_opts"`
	PluginArgs    string `json:"plugin_args"`
	Remarks       string `json:"remarks"`
	Group         string `json:"group"`
	Timeout       int    `json:"timeout"`
	LocalPort     int    `json:"local_port"`
	LocalAddress  string `json:"local_address"`
	Fastopen      bool   `json:"fastopen"`
	Obfsparam     string `json:"obfs_param"`
}

type Result struct {
	Status  string       `json:"status"`
	Current ShadowConfig `json:"current"`
	Result  []ResultInfo `json:"results"`
}

type ResultInfo struct {
	Group   string  `json:"group"`
	Remarks string  `json:"remarks"`
	Loss    float64 `json:"loss"`
	Ping    float64 `json:"ping"`
	GPing   float64 `json:"gPing"`
}

type StartConfigs struct {
	TestMethod   string              `json:"testMethod"`
	TestMode     string              `json:"testMode"`
	Colors       string              `json:"colors"`
	SortMethod   string              `json:"sortMethod"`
	UseSSRCSharp bool                `json:"useSsrcSharp"`
	Group        string              `json:"group"`
	Configs      []*SubscriptionResp `json:"configs"`
}

func NewStartConfigs(testMethod string, testMode string, configs []*SubscriptionResp) *StartConfigs {
	return &StartConfigs{
		TestMethod: testMethod,
		TestMode:   testMode,
		Configs:    configs,
	}
}
