package reports

import (
	"strings"

	"github.com/golovers/kiki/backend/issues"

	"github.com/golovers/kiki/backend/types"
	"github.com/golovers/kiki/backend/utils"
)

type reportSvc struct {
}

//NewSimReport return a simple report service
func NewSimReport() Reports {
	return &reportSvc{}
}

type FilterFunc func(issue *types.Issue) bool
type AggrFunc func(issue *types.Issue) float64

var defectFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "defect", "bug")
}

var storyFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "story", "improvement", "enhancement")
}

var teamSprintStoryFilter = func(team string, sprint string) FilterFunc {
	return func(issue *types.Issue) bool {
		if !storyFilter(issue) {
			return false
		}
		if team != "*" && !utils.OneOf(team, issue.Labels...) {
			return false
		}
		if sprint != "*" && strings.ToLower(sprint) != strings.ToLower(issue.Sprint) {
			return false
		}
		return true
	}
}

var epicFilter = func(epic string) FilterFunc {
	return func(issue *types.Issue) bool {
		if epic != "*" && epic != issue.EpicLink {
			return false
		}
		return true
	}
}

var multipleFilters = func(filters ...FilterFunc) FilterFunc {
	return func(issue *types.Issue) bool {
		for _, f := range filters {
			ok := f(issue)
			if !ok {
				return false
			}
		}
		return true
	}
}

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
	rs := make(TeamSprintStatus)
	rs["_all"] = SprintStatus{
		Defects: svc.Defects(teamSprintStoryFilter("*", sprint)),
		Stories: svc.Stories(teamSprintStoryFilter("*", sprint)),
	}
	for _, team := range teams {
		rs[team] = SprintStatus{
			Defects: svc.Defects(teamSprintStoryFilter(team, sprint)),
			Stories: svc.Stories(teamSprintStoryFilter(team, sprint)),
		}
	}
	otherTeamsFilter := func(issue *types.Issue) bool {
		for _, team := range teams {
			if teamSprintStoryFilter(team, sprint)(issue) {
				return false
			}
		}
		return true
	}
	rs["other"] = SprintStatus{
		Defects: svc.Defects(otherTeamsFilter),
		Stories: svc.Stories(otherTeamsFilter),
	}
	return rs
}
