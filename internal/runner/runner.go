package runner

import "sync/atomic"

type Runner struct {
	status int32  // Status store runner status at local
	Name   string `json:"name"`
	Host   *Host  `json:"host"`
	Admins []int  `json:"admins"`
}

type Host struct {
	Address string `json:"address"`
	Key     string `json:"key"`
}

func (h *Host) GetURL() string {
	return h.Address
}

const (
	Pending = iota // 0 == Pending
	Working        // 1 == Working
)

// IsPending return boolen value about runner is pending or not
func (r *Runner) IsPending() bool {
	return r.GetRunnerStatus() == Pending
}

// IsPending return boolen value about runner is working or not
func (r *Runner) IsWorking() bool {
	return r.GetRunnerStatus() == Working
}

// GetRunnerStatus return current status
func (r *Runner) GetRunnerStatus() int32 {
	return atomic.LoadInt32(&r.status)
}

// HangUp changed runner status to pending
func (r *Runner) HangUp() {
	atomic.CompareAndSwapInt32(&r.status, Working, Pending)
}

// Activate changed runner status to working
func (r *Runner) Activate() {
	atomic.CompareAndSwapInt32(&r.status, Pending, Working)
}
