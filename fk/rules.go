package fk

import (
	"encoding/json"
	"time"

	"github.com/elgs/jsonql"
)

type Rules []Ruler

type Ruler interface {
	Run()
	SchedSpec() string
	Name() string
	ID() string
}

type Rule struct {
	ID             string
	Name           string
	Query          string
	Schedule       string // when to run this rule
	EventThreshold int    // when to fire event, this can be different between alarm types
	AlarmThreshold int    // how many continous event to trigger alarm
}

// Event represent an even fired by rule engine
type Event struct {
	ID        string
	Timestamp time.Time
	Violated  bool
	Issues    Issues
}

// Run the issues list and return violated list
func (ca *Rule) Run(issues Issues) (*Event, error) {
	jsonBytes, err := json.Marshal(issues)
	if err != nil {
		return &Event{}, err
	}
	query, err := jsonql.NewStringQuery(string(jsonBytes))
	if err != nil {
		return &Event{}, err
	}
	val, err := query.Query(ca.Query)
	if err != nil {
		return &Event{}, err
	}
	var rissues Issues
	d, _ := json.Marshal(val)
	json.Unmarshal(d, &rissues)
	evn := &Event{
		ID:        GenID(),
		Timestamp: time.Now(),
		Violated:  len(rissues) >= ca.EventThreshold, //TODO fix me
		Issues:    rissues,
	}
	return evn, nil
}
