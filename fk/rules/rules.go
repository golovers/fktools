package rules

var svc RuleSvc

type RuleSvc interface {
	Load() ([]*Rule, error)
}

func SetRuleSvc(s RuleSvc) {
	svc = s
}

func Load() ([]*Rule, error) {
	return svc.Load()
}
