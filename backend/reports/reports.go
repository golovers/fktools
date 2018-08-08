package reports

import (
	"strings"

	"github.com/golovers/kiki/backend/types"
	"github.com/golovers/kiki/backend/utils"
)

// PriorityAndStatus priority and status in two dementions values
type PriorityAndStatus struct {
	Critical *StatusSummary
	Major    *StatusSummary
	Minor    *StatusSummary
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

type reportSvc struct {
	issues types.Issues
}

type filterFunc func(issue *types.Issue) bool
type aggrFunc func(issue *types.Issue) float64

var defectFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "defect", "bug")
}

var storyFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "story", "improvement", "enhancement")
}

var aggrCount = func(issue *types.Issue) float64 {
	return 1.0
}
var aggrStoryPoints = func(issue *types.Issue) float64 {
	return issue.StoryPoints
}

func (svc *reportSvc) Status(filter filterFunc, aggr aggrFunc) *PriorityAndStatus {
	status := &PriorityAndStatus{}
	for _, issue := range svc.issues {
		if filter(issue) {
			status.update(issue.Priority, issue.Status, aggr(issue))
		}
	}
	return status
}

func (svc *reportSvc) DefectCountStatus() *PriorityAndStatus {
	return svc.Status(defectFilter, aggrCount)
}

func (svc *reportSvc) StoryPriorityCountStatus() *PriorityAndStatus {
	return svc.Status(storyFilter, aggrCount)
}

func (svc *reportSvc) StoryPriorityStoryPointStatus() *PriorityAndStatus {
	return svc.Status(storyFilter, aggrStoryPoints)
}
