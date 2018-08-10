package reports

import (
	"github.com/golovers/kiki/backend/issues"

	"github.com/golovers/kiki/backend/types"
)

type reportSvc struct {
}

//NewSimReport return a simple report service
func NewSimReport() Reports {
	return &reportSvc{}
}

type AggrFunc func(issue *types.Issue) float64

var aggrCount = func(issue *types.Issue) float64 {
	return 1.0
}
var aggrStoryPoints = func(issue *types.Issue) float64 {
	return issue.StoryPoints
}

func (svc *reportSvc) Issues() types.Issues {
	issues, err := issues.Load()
	if err != nil {
		return types.Issues{}
	}
	return issues
}

func (svc *reportSvc) status(aggr AggrFunc, filters ...FilterFunc) *PriorityAndStatus {
	status := &PriorityAndStatus{
		Critical: &StatusSummary{},
		Major:    &StatusSummary{},
		Minor:    &StatusSummary{},
	}
	cbFilters := multipleFilters(filters...)
	for _, issue := range svc.Issues() {
		if cbFilters(issue) {
			status.update(issue.Priority, issue.Status, aggr(issue))
		}
	}
	return status
}

func (svc *reportSvc) Defects(filters ...FilterFunc) *PriorityAndStatus {
	filters = append(filters, defectFilter)
	return svc.status(aggrCount, filters...)
}

func (svc *reportSvc) Stories(filters ...FilterFunc) *StoryStatus {
	filters = append(filters, storyFilter)
	points := svc.status(aggrStoryPoints, filters...).Overview()
	count := svc.status(aggrCount, filters...).Overview()

	return &StoryStatus{
		Count:  count,
		Points: points,
	}
}

func (svc *reportSvc) Sprint(sprint string, teams ...string) TeamSprintStatus {
	rs := TeamSprintStatus{
		Sprint:     sprint,
		TeamStatus: make(map[string]SprintStatus),
	}
	rs.TeamStatus["z_total_z"] = SprintStatus{
		Defects: svc.Defects(teamSprintDefectFilter(sprint, "*")),
		Stories: svc.Stories(teamSprintStoryFilter(sprint, "*")),
	}
	for _, team := range teams {
		rs.TeamStatus[team] = SprintStatus{
			Defects: svc.Defects(teamSprintDefectFilter(sprint, team)),
			Stories: svc.Stories(teamSprintStoryFilter(sprint, team)),
		}
	}

	rs.TeamStatus["z_other_z"] = SprintStatus{
		Defects: svc.Defects(otherTeamsSprintDefectFilter(sprint, teams)),
		Stories: svc.Stories(otherTeamsSprintStoryFilter(sprint, teams)),
	}
	return rs
}
