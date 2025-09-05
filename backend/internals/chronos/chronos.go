package chronos

import (
	"time"
)

type Timer struct {
	isActive  bool
	startLap  time.Time
	spentTime time.Duration
	EndTime   time.Duration
}

func NewTimer(end time.Duration) *Timer {
	return &Timer{
		spentTime: 0,
		EndTime:   end,
		isActive:  false,
	}
}

func (t *Timer) Start() {
	t.startLap = time.Now()
	t.isActive = true
}

func (t *Timer) Stop() {
	t.spentTime += time.Since(t.startLap)
	t.isActive = false
}

func (t *Timer) Spent() time.Duration {
	if t.isActive {
		t.Stop()
		t.Start()
	}
	return t.spentTime
}

func (t *Timer) SpentInt() int {
	return int(t.Spent().Seconds())
}
