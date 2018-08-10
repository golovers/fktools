package reports

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type Reports interface {
	Defects(filters ...FilterFunc) *PriorityAndStatus
	Stories(filters ...FilterFunc) *StoryStatus
	Sprint(sprint string, teams ...string) TeamSprintStatus
}

type TeamSprintStatus struct {
	Sprint     string
	TeamStatus map[string]SprintStatus
}

type SprintStatus struct {
	Defects *PriorityAndStatus
	Stories *StoryStatus
}

// PriorityAndStatus priority and status in two dementions values
type PriorityAndStatus struct {
	Critical *StatusSummary
	Major    *StatusSummary
	Minor    *StatusSummary
}

type StoryStatus struct {
	Count  *StatusSummary
	Points *StatusSummary
}

// Overview return Overview status by summing values of the 2 dimentions values base on status of the issue
func (ps *PriorityAndStatus) Overview() *StatusSummary {
	return ps.Critical.Merge(ps.Major).Merge(ps.Minor)
}

func (ps *PriorityAndStatus) update(pri string, status string, v float64) {
	switch strings.ToLower(pri) {
	case "critical":
		ps.Critical.update(status, v)
	case "major":
		ps.Major.update(status, v)
	case "minor":
		ps.Minor.update(status, v)
	default:
		logrus.Errorf("not supported priority: %s", pri)
	}
}

// StatusSummary hold over all status of issues by counting or summing story points
type StatusSummary struct {
	Open       float64
	InProgress float64
	Resolved   float64
	Reopened   float64
	Closed     float64
}

// Total return sum of all values
func (ov *StatusSummary) Total() float64 {
	return ov.Open + ov.InProgress + ov.Reopened + ov.Resolved + ov.Closed
}

func (ov *StatusSummary) update(status string, v float64) {
	switch status {
	case "open":
		ov.Open += v
	case "inprogress":
		ov.InProgress += v
	case "resolved":
		ov.Resolved += v
	case "reopened":
		ov.Resolved += v
	case "closed":
		ov.Closed += v
	default:
		logrus.Infof("not supported status: %s", status)
	}
}

// TotalOpen return sum of all kind status diffrent from close
func (ov *StatusSummary) TotalOpen() float64 {
	return ov.Open + ov.InProgress + ov.Reopened + ov.Resolved
}

// Merge the current with the given status into one
func (ov *StatusSummary) Merge(o *StatusSummary) *StatusSummary {
	return &StatusSummary{
		Open:       ov.Open + o.Open,
		InProgress: ov.InProgress + o.InProgress,
		Reopened:   ov.Reopened + o.Reopened,
		Resolved:   ov.Resolved + o.Resolved,
		Closed:     ov.Closed + o.Closed,
	}
}

var svc Reports

// SetReports set report service
func SetReports(s Reports) {
	svc = s
}

// Defects defect status by priority and status
func Defects() *PriorityAndStatus {
	return svc.Defects()
}

// Stories story status by priority and status
func Stories() *StoryStatus {
	return svc.Stories()
}

func Sprint(sprint string, teams ...string) TeamSprintStatus {
	return svc.Sprint(sprint, teams...)
}
