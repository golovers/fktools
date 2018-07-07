package rules

import (
	"encoding/json"
	"github.com/golovers/fktools/fk/db"
)

var keyRules = []byte("rules")

type simSvc struct{}

func NewSimRuleSvc() RuleSvc {
	return &simSvc{}
}

// AllIssues load all existing rules in db
func (s *simSvc) Load() ([]*Rule, error) {
	data, err := db.Get(keyRules)
	if err != nil {
		return []*Rule{}, err
	}
	var rules []*Rule
	err = json.Unmarshal(data, &rules)
	if err != nil {
		return []*Rule{}, err
	}
	return rules, nil
}
