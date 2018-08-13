package types

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/golovers/kiki/backend/utils"
	"github.com/sirupsen/logrus"
)

var attrs = []string{"Issue Type", "Key", "Status", "Summary", "Priority", "Story Points", "Epic", "Sprint",
	"Fix Version/svc", "Reporter", "Affected Version/svc", "Assignee", "Components", "Created", "Labels",
	"Resolution", "Resolved", "Time Original Estimate", "Time Estimate",
	"Time Spent"}

type Issues []*Issue

// Issue represent an issue
type Issue struct {
	IssueType            string
	Key                  string
	Status               string
	Summary              string
	Priority             string
	StoryPoints          float64
	EpicLink             string
	Sprint               string
	FixVersions          []string
	Reporter             string
	AffectsVersions      []string
	Assignee             string
	Components           []string
	Created              *time.Time
	Labels               []string
	Resolution           string
	Resolved             *time.Time
	TimeOriginalEstimate int
	TimeEstimate         int
	TimeSpent            int
}

// ToCSV return csv string of the issue
func (issue *Issue) ToCSV() string {
	fmtStr := ""
	for i := 0; i < len(attrs); i++ {
		fmtStr += "\"%v\""
		if i < len(attrs)-1 {
			fmtStr += ","
		}
	}
	fmtStr += "\n"
	str := fmt.Sprintf(fmtStr,
		issue.IssueType, issue.Key, issue.Status, utils.EscapeSpecialChars(issue.Summary), issue.Priority, issue.StoryPoints,
		issue.EpicLink, issue.Sprint, utils.StringsToString(issue.FixVersions), issue.Reporter, utils.StringsToString(issue.AffectsVersions),
		issue.Assignee, utils.StringsToString(issue.Components), issue.Created,
		utils.EscapeSpecialChars(utils.StringsToString(issue.Labels)), issue.Resolution, issue.Resolved,
		issue.TimeOriginalEstimate, issue.TimeEstimate, issue.TimeSpent)
	return str
}

// ToCSV export issues into csv file
func (issues *Issues) ToCSV(file string) (string, error) {
	f, err := os.Create(file)
	if err != nil {
		return "", err
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%svc\n", strings.Join(attrs, ",")))
	for i, issue := range *issues {
		f.WriteString(issue.ToCSV())
		if i%100 == 0 {
			f.Sync()
		}
	}
	logrus.Infof("exported issues to: %svc", file)
	return f.Name(), nil
}

// IssuesToJSONString convert list of issues to a json string
func IssuesToJSONString(issues Issues) (string, error) {
	data, err := json.Marshal(issues)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (iss Issues) Sort() {
	sort.Slice(iss, func(i, j int) bool {
		return strings.Compare(iss[i].Status, iss[j].Status) == 0
	})
}
