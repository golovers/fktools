package reports

import (
	"fmt"
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

func (svc *reportSvc) Issues() types.Issues {
	issues, err := issues.Load()
	if err != nil {
		return types.Issues{}
	}
	return issues
}

func (svc *reportSvc) status(aggr AggrFunc, filters ...FilterFunc) *PriorityAndStatus {
	status, _ := svc.statusAndIssues(aggr, filters...)
	return status
}

func (svc *reportSvc) statusAndIssues(aggr AggrFunc, filters ...FilterFunc) (*PriorityAndStatus, types.Issues) {
	status := &PriorityAndStatus{
		Critical: &StatusSummary{},
		Major:    &StatusSummary{},
		Minor:    &StatusSummary{},
	}
	filIssues := make(types.Issues, 0)
	cbFilters := multipleFilters(filters...)
	for _, issue := range svc.Issues() {
		if cbFilters(issue) {
			filIssues = append(filIssues, issue)
			status.update(issue.Priority, issue.Status, aggr(issue))
		}
	}
	return status, filIssues
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

func (svc *reportSvc) Sprint(sprint string, labels ...string) TeamSprintStatus {
	rs := TeamSprintStatus{
		Sprint:     sprint,
		TeamStatus: make(map[string]SprintStatus),
	}
	rs.TeamStatus["z_total_z"] = SprintStatus{
		Defects: svc.Defects(labelSprintDefectFilter(sprint, labels...)),
		Stories: svc.Stories(labelSprintStoryFilter(sprint, labels...)),
	}
	if len(labels) > 0 && labels[0] == "" {
		return rs
	}
	for _, label := range labels {
		rs.TeamStatus[label] = SprintStatus{
			Defects: svc.Defects(labelSprintDefectFilter(sprint, label)),
			Stories: svc.Stories(labelSprintStoryFilter(sprint, label)),
		}
	}
	if utils.OneOf("_other_", labels...) {
		rs.TeamStatus["z_other_z"] = SprintStatus{
			Defects: svc.Defects(otherLabelsSprintDefectFilter(sprint, labels)),
			Stories: svc.Stories(otherLabelsSprintStoryFilter(sprint, labels)),
		}
	}
	return rs
}

type GroupStatus struct {
	Name          string
	StoriesStatus *StatusSummary
	DefectsStatus *PriorityAndStatus
	Stories       types.Issues
	Defects       types.Issues
}

func (svc *reportSvc) groupStatus(filters ...FilterFunc) *GroupStatus {
	status := new(GroupStatus)
	storyFilters := append(filters, storyFilter)
	storyStatus, stories := svc.statusAndIssues(aggrStoryPoints, multipleFilters(storyFilters...))
	status.StoriesStatus = storyStatus.Overview()
	status.Stories = stories
	defectFilters := append(filters, defectFilter)
	status.DefectsStatus, status.Defects = svc.statusAndIssues(aggrCount, defectFilters...)
	return status
}

func (svc *reportSvc) EpicStatus(epic string, fixVersions []string, labels []string, sprint string) *GroupStatus {
	status := svc.groupStatus(epicFilter(epic), fixVersionsFilter(fixVersions...), oneOfTheseLabels(labels...), sprintFilter(sprint))
	status.Name = fmt.Sprintf("Epic: %s - FixVersions: %s - Labels: %s - Sprint: %s", epic, strings.Join(fixVersions, ", "), strings.Join(labels, ", "), sprint)
	return status
}
