package fk

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	plug     Plugin
	tx       Transformer
	db       Database
	sched    Scheduler
	keyIssue = []byte("issues")
	keyRules = []byte("rules")
)

// SetPlugin set the plugin to JIRA, VersionOne, Phabricator,...
func SetPlugin(p Plugin) {
	plug = p
}

// SetDatabase set database to be used
func SetDatabase(database Database) {
	db = database
}

// SetTransformer set the StatusTransformer to be used
func SetTransformer(s Transformer) {
	tx = s
}

// SetScheduler set a scheduler to be used
func SetScheduler(sch Scheduler) {
	sched = sch
}

// AllIssues return all issues
func AllIssues() (Issues, error) {
	data, err := db.Get(keyIssue)
	if err != nil {
		return Issues{}, err
	}
	var issues Issues
	err = json.Unmarshal(data, &issues)
	if err != nil {
		return Issues{}, err
	}
	transform(issues)
	return issues, nil
}

// Sync sync data from remote to local for later use
func Sync() {
	logrus.Info("syncing data from remote...")
	issues, err := plug.AllIssues()
	if err != nil {
		logrus.Error("failed to sync data from remote", err)
		//TODO send mail to admin
		return
	}
	saveIssues(issues)
	logrus.Info("finished sync data from remote")
}

func saveIssues(issues Issues) {
	data, err := json.Marshal(issues)
	if err != nil {
		logrus.Errorf("failed to save issues: %v", err)
	}
	db.Put(keyIssue, data)
}

func transform(issues Issues) {
	if tx == nil {
		tx = &CommonTransformer{}
	}
	for _, issue := range issues {
		issue = tx.Transform(*issue)
	}
}

// ToCSV export all issues to CSV file
func ToCSV(file string) (string, error) {
	issues, err := AllIssues()
	if err != nil {
		return file, err
	}
	return issues.ToCSV(file)
}

// ToCSV export issues into csv file
func (issues *Issues) ToCSV(file string) (string, error) {
	f, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s\n", strings.Join(attrs, ",")))
	for i, issue := range *issues {
		f.WriteString(issue.ToCSV())
		if i%100 == 0 {
			f.Sync()
		}
	}
	logrus.Infof("exported issues to: %s", file)
	return f.Name(), nil
}

// SchedSync sync task
func SchedSync(spec string) {
	sched.Sched(spec, Sync)
}

// LoadRules load all existing rules in db
func LoadRules() ([]*FKRule, error) {
	data, err := db.Get(keyRules)
	if err != nil {
		return []*FKRule{}, err
	}
	var rules []*FKRule
	err = json.Unmarshal(data, &rules)
	if err != nil {
		return []*FKRule{}, err
	}
	return rules, nil
}

// SchedRules schedule running rules
func SchedRules() {
	rules, err := LoadRules()
	if err != nil {
		logrus.Errorf("failed to load rules: %v", err)
	}
	for _, r := range rules {
		logrus.Infof("scheduled for rule \"%s-%s\"", r.Name, r.ID)
		sched.Sched(r.Schedule, func() {
			evn, err := r.Run(AllIssues)
			logrus.Infof("raised event %s, violated: %v, # issues: %d", evn.ID, evn.Violated, len(evn.Issues))
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

func isReachAlarmThreshold(r *FKRule, events []*Event) bool {
	count := 0
	for _, e := range events {
		if e.Violated {
			count++
		}
	}
	return count >= r.AlarmThreshold
}

func fireAlarm(r *FKRule, events []*Event) {
	//TODO implement me
	logrus.Infof("firing alarm for rule %s-%s: %d continuos time violated", r.Name, r.ID, r.AlarmThreshold)
}

// saveEvent save the given event and return the last n event from db
func saveEvent(r *FKRule, evn *Event) ([]*Event, error) {
	eventData, err := db.Get([]byte(r.ID))
	if err != nil || len(eventData) == 0 {
		// not yet has any data, just update this into the db
		events := make([]*Event, 0)
		events = append(events, evn)
		jsonStr, _ := json.Marshal(events)
		db.Put([]byte(r.ID), jsonStr)
	}
	var events []*Event
	err = json.Unmarshal(eventData, &events)
	if err != nil {
		return []*Event{}, fmt.Errorf("failed to unmarshall events data of rule %s-%s, error: %v", r.Name, r.ID, err)
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
