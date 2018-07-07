package sched

// Scheduler wraps all functions that a scheduler must has
type Scheduler interface {
	Start()
	Schedule(spec string, f func()) error
	Stop()
}

var sched Scheduler

// SetScheduler set a scheduler to be used
func SetScheduler(sch Scheduler) {
	sched = sch
}

func Start() {
	sched.Start()
}

func Schedule(spec string, f func()) error {
	return sched.Schedule(spec, f)
}

func Stop() {
	sched.Stop()
}
