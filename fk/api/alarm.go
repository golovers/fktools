package api

import (
	"encoding/json"
	"fmt"

	"github.com/golovers/fktools/fk/db"
	"github.com/golovers/fktools/fk/iss"
	"github.com/golovers/fktools/fk/rules"
	"github.com/golovers/fktools/fk/sched"
	"github.com/sirupsen/logrus"
)

// SchedRules schedule running rules
func SchedRules() {
	rules, err := rules.Load()
	if err != nil {
		logrus.Errorf("failed to load rules: %v", err)
	}
	for _, r := range rules {
		logrus.Infof("scheduled for rule \"%s-%s\"", r.Name, r.ID)
		sched.Schedule(r.Schedule, func() {
			evn, err := r.Run()
			logrus.Infof("raised event %s, violated: %v, # issue: %d", evn.ID, evn.Violated, len(evn.Issues))
			if err != nil {
				logrus.Errorf("failed to run the rule %s-%s, error: %v", r.Name, r.ID, err)
			}
			events, err := saveEvent(r, evn)
			if err != nil {
				logrus.Error(err)
			}
			if isReachAlarmThreshold(r, events) {
				fireAlarm(r, events)
			}
		})
	}
}

func isReachAlarmThreshold(r *rules.Rule, events []*rules.Event) bool {
	count := 0
	for _, e := range events {
		if e.Violated {
			count++
		}
	}
	return count >= r.AlarmThreshold
}

func fireAlarm(r *rules.Rule, events []*rules.Event) {
	//TODO implement me
	logrus.Infof("firing alarm for rule %s-%s: %d continuos time violated", r.Name, r.ID, r.AlarmThreshold)
}

// saveEvent save the given event and return the last n event from db
func saveEvent(r *rules.Rule, evn *rules.Event) ([]*rules.Event, error) {
	eventData, err := db.Get([]byte(r.ID))
	if err != nil || len(eventData) == 0 {
		// not yet has any data, just update this into the db
		events := make([]*rules.Event, 0)
		events = append(events, evn)
		jsonStr, _ := json.Marshal(events)
		db.Put([]byte(r.ID), jsonStr)
	}
	var events []*rules.Event
	err = json.Unmarshal(eventData, &events)
	if err != nil {
		return []*rules.Event{}, fmt.Errorf("failed to unmarshall events data of rule %s-%s, error: %v", r.Name, r.ID, err)
	}
	events = append(events, evn)
	jsonStr, _ := json.Marshal(events)
	db.Put([]byte(r.ID), jsonStr)

	thresholdIdx := r.AlarmThreshold
	if r.AlarmThreshold-thresholdIdx < 0 {
		thresholdIdx = 0
	}
	return events[thresholdIdx:], nil
}

// ToCSV export all iss to CSV file
func ToCSV(file string) (string, error) {
	issues, err := iss.Load()
	if err != nil {
		return file, err
	}
	return issues.ToCSV(file)
}

// SchedSync sync task
func SchedSync(spec string) {
	sched.Schedule(spec, iss.Sync)
}
