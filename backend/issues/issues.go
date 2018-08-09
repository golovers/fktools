package issues

import "github.com/golovers/kiki/backend/types"
import "github.com/golovers/kiki/backend/sched"

var svc IssueSvc

type IssueSvc interface {
	Load() (types.Issues, error)
	Sync()
}

func SetIssueSvc(s IssueSvc) {
	svc = s
}

func Load() (types.Issues, error) {
	return svc.Load()
}

func Sync() {
	svc.Sync()
}

func SchedSync(spec string) {
	sched.Schedule(spec, Sync)
}
