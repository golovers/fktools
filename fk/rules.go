package fk

import (
	"encoding/json"
	"time"

	"github.com/elgs/jsonql"
)

type Rules []Ruler
type IssuesFunc func() (Issues, error)

type Ruler interface {
	Run(IssuesFunc) (*Event, error)
	Sched() string
	Name() string
	ID() string
	AlarmThreshold() int
}

// Event represent an even fired by rule engine
type Event struct {
	ID        string
	Timestamp time.Time
	Violated  bool
	Issues    Issues
}

type FKRule struct {
	ID             string
	Name           string
	Query          string
	Schedule       string // when to run this rule
	EventThreshold int    // when to fire event, this can be different between alarm types
	AlarmThreshold int    // how many continous event to trigger alarm
}

// Run the issues list and return violated list
func (cex *FKRule) Run(f IssuesFunc) (*Event, error) {
	issues, err := f()
	if err != nil {
		return &Event{}, err
	}
	jsonBytes, err := json.Marshal(issues)
	query, err := jsonql.NewStringQuery(string(jsonBytes))
	if err != nil {
		return &Event{}, err
	}
	val, err := query.Query(cex.Query)
	if err != nil {
		return &Event{}, err
	}
	var rissues Issues
	d, _ := json.Marshal(val)
	json.Unmarshal(d, &rissues)
	evn := &Event{
		ID:        GenID(),
		Timestamp: time.Now(),
		Violated:  len(rissues) >= cex.EventThreshold,
		Issues:    rissues,
	}
	return evn, nil
}
