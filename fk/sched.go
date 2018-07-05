package fk

import "github.com/robfig/cron"

// Scheduler wraps all functions that a scheduler must has
type Scheduler interface {
	Start()
	Sched(spec string, f func()) error
	Stop()
}

type cronScheduler struct {
	sched *cron.Cron
}

// NewCronScheduler return a simple cron scheduler
func NewCronScheduler() Scheduler {
	return &cronScheduler{
		sched: cron.New(),
	}
}

func (cs *cronScheduler) Start() {
	cs.sched.Start()
}

func (cs *cronScheduler) Stop() {
	cs.sched.Stop()
}

func (cs *cronScheduler) Sched(spec string, f func()) error {
	return cs.sched.AddFunc(spec, f)
}
