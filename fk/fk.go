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
func LoadRules() ([]Ruler, error) {
	data, err := db.Get(keyRules)
	if err != nil {
		return []Ruler{}, err
	}
	var rules []Ruler
	err = json.Unmarshal(data, &rules)
	if err != nil {
		return []Ruler{}, err
	}
	return rules, nil
}

// SchedRules schedule running rules
func SchedRules(spec string) {

}
