package rules

import (
	"encoding/json"
	"github.com/elgs/jsonql"
	"github.com/golovers/fktools/fk/iss"
	"github.com/golovers/fktools/fk/types"
	"github.com/golovers/fktools/fk/utils"
	"time"
)

// Event represent an even fired by rule engine
type Event struct {
	ID        string
	Timestamp time.Time
	Violated  bool
	Issues    types.Issues
}

type Rule struct {
	ID             string
	Name           string
	Query          string
	Schedule       string // when to run this rule
	EventThreshold int    // when to fire event, this can be different between alarm types
	AlarmThreshold int    // how many continous event to trigger alarm
}

// Run the iss list and return violated list
func (cex *Rule) Run() (*Event, error) {
	issues, err := iss.Load()
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
	var rissues types.Issues
	d, _ := json.Marshal(val)
	json.Unmarshal(d, &rissues)
	evn := &Event{
		ID:        utils.GenID(),
		Timestamp: time.Now(),
		Violated:  len(rissues) >= cex.EventThreshold,
		Issues:    rissues,
	}
	return evn, nil
}
