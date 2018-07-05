package fk

import (
	"fmt"
	"strings"
	"time"
)

var attrs = []string{"Issue Type", "Key", "Status", "Summary", "Priority", "Story Points", "Epic", "Sprint",
	"Fix Version/s", "Reporter", "Affected Version/s", "Assignee", "Components", "Created", "Labels",
	"Resolution", "Resolved", "Time Original Estimate", "Time Estimate",
	"Time Spent"}

func StringToStrings(s string) []string {
	return strings.Split(s, ",")
}

func StringsToString(s []string) string {
	return strings.Join(s, ",")
}

func EscapeSpecialChars(s string) string {
	return strings.Replace(s, "\"", "\\\"", -1)
}

type Issues []*Issue

type Plugin interface {
	AllIssues() (Issues, error)
}

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
		issue.IssueType, issue.Key, issue.Status, EscapeSpecialChars(issue.Summary), issue.Priority, issue.StoryPoints,
		issue.EpicLink, issue.Sprint, StringsToString(issue.FixVersions), issue.Reporter, StringsToString(issue.AffectsVersions),
		issue.Assignee, StringsToString(issue.Components), issue.Created,
		EscapeSpecialChars(StringsToString(issue.Labels)), issue.Resolution, issue.Resolved,
		issue.TimeOriginalEstimate, issue.TimeEstimate, issue.TimeSpent)
	return str
}

type Transformer interface {
	Transform(issue Issue) *Issue
}

type CommonTransformer struct {
}

func (tx *CommonTransformer) status(s string) string {
	if contains([]string{"Critical", "High", "Highest", "Blocker"}, s) {
		return "Critical"
	}
	if contains([]string{"Major"}, s) {
		return "Major"
	}
	if contains([]string{"Low", "Minor", "Medium", ""}, s) {
		return "Minor"
	}
	return s
}

func (tx *CommonTransformer) Transform(issue Issue) *Issue {
	newIssue := &issue
	newIssue.Status = tx.status(newIssue.Status)
	if newIssue.Priority == "" {
		newIssue.Priority = " "
	}
	return newIssue
}
