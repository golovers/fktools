package issues

import "github.com/golovers/kiki/backend/types"

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
