package utils

import "time"

type PeriodicFunc struct {
	task        func()
	ticker      *time.Ticker
	channelDone chan bool
}

func CreatePeriodic(interval time.Duration, task func()) *PeriodicFunc {
	return &PeriodicFunc{
		task:        task,
		ticker:      time.NewTicker(interval),
		channelDone: make(chan bool),
	}
}

func (periodic *PeriodicFunc) Start() {
	go periodic.wrap()
}

func (periodic *PeriodicFunc) Stop() {
	periodic.channelDone <- true
}

func (periodic *PeriodicFunc) wrap() {
	for {
		select {
		case <-periodic.channelDone:
			periodic.ticker.Stop()
			return
		case <-periodic.ticker.C:
			periodic.task()
		}
	}
}
