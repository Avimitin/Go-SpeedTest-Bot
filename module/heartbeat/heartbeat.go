package heartbeat

import "time"

type HeartBeat struct {
	period time.Duration
	ticker *time.Ticker
}

// NewHeartBeat return a new HeartBeat that containing a time.Ticker
// and user specific period.
// It is used for easily reset the ticker.
func NewHeartBeat(interval int) *HeartBeat {
	p := (time.Duration(interval) * time.Second) / 2
	return &HeartBeat{
		period: p,
		ticker: time.NewTicker(p),
	}
}

// C return the ticker channel
func (hb *HeartBeat) C() <-chan time.Time {
	return hb.ticker.C
}

// Stop turns off a ticker. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a concurrent goroutine
// reading from the channel from seeing an erroneous "tick".
func (hb *HeartBeat) Stop() {
	hb.ticker.Stop()
}

// Reset stops a ticker and resets its period to the specified duration.
// The next tick will arrive after the new period elapses.
func (hb *HeartBeat) Reset() {
	hb.ticker.Reset(hb.period)
}
