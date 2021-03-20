package runner

import "sync/atomic"

type Runner struct {
	// status store runner status at local
	status int32

	// c is a channel to connect with runner inner goroutine,
	// should use Runner.NewChan to initialize it and get it,
	// and use Runner.CloseChan to close it.
	c chan int32

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

func (r *Runner) NewChan() chan int32 {
	r.c = make(chan int32)
	return r.c
}

func (r *Runner) CloseChan() {
	close(r.c)
}

// IsPending return boolean value about runner is pending or not
func (r *Runner) IsPending() bool {
	return r.GetRunnerStatus() == Pending
}

// IsPending return boolean value about runner is working or not
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

// HasAccess testify the given have permission to use this
// runner or not.
func (r *Runner) HasAccess(id int) bool {
	for _, admin := range r.Admins {
		if id == admin {
			return true
		}
	}
	return false
}
