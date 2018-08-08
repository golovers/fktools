package sched

import "github.com/robfig/cron"

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

func (cs *cronScheduler) Schedule(spec string, f func()) error {
	return cs.sched.AddFunc(spec, f)
}
