package reports

import (
	"strings"

	"github.com/golovers/kiki/backend/types"
	"github.com/golovers/kiki/backend/utils"
)

// FilterFunc check if the issue match a certain conditions
type FilterFunc func(issue *types.Issue) bool

var reverseFilter = func(filter FilterFunc) FilterFunc {
	return func(issue *types.Issue) bool {
		return !filter(issue)
	}
}

var defectFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "defect", "bug")
}

var storyFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "story", "improvement", "enhancement")
}

var oneOfTheseLabels = func(labels ...string) FilterFunc {
	return func(issue *types.Issue) bool {
		return utils.AnyOf(labels, issue.Labels...)
	}
}

var notTheseLabelsFilter = func(labels ...string) FilterFunc {
	return reverseFilter(oneOfTheseLabels(labels...))
}

var sprintFilter = func(sprint string) FilterFunc {
	return func(issue *types.Issue) bool {
		return sprint == "*" || utils.OneOf(sprint, issue.Sprint)
	}
}

var labelSprintFilter = func(issueTypeFilter FilterFunc, sprint string, label ...string) FilterFunc {
	return multipleFilters(issueTypeFilter, sprintFilter(sprint), oneOfTheseLabels(label...))
}

var labelSprintStoryFilter = func(sprint string, labels ...string) FilterFunc {
	return labelSprintFilter(storyFilter, sprint, labels...)
}

var labelSprintDefectFilter = func(sprint string, labels ...string) FilterFunc {
	return labelSprintFilter(defectFilter, sprint, labels...)
}

var otherLabelsSprintFilter = func(issueTypeFilter FilterFunc, sprint string, labels []string) FilterFunc {
	return multipleFilters(issueTypeFilter, sprintFilter(sprint), notTheseLabelsFilter(labels...))
}

var otherLabelsSprintDefectFilter = func(sprint string, labels []string) FilterFunc {
	return otherLabelsSprintFilter(defectFilter, sprint, labels)
}

var otherLabelsSprintStoryFilter = func(sprint string, labels []string) FilterFunc {
	return otherLabelsSprintFilter(storyFilter, sprint, labels)
}

var epicFilter = func(epic string) FilterFunc {
	return func(issue *types.Issue) bool {
		return epic == "*" || utils.OneOf(epic, issue.EpicLink)
	}
}

var multipleFilters = func(filters ...FilterFunc) FilterFunc {
	return func(issue *types.Issue) bool {
		for _, f := range filters {
			if ok := f(issue); !ok {
				return false
			}
		}
		return true
	}
}

var fixVersionsFilter = func(versions ...string) FilterFunc {
	return func(issue *types.Issue) bool {
		return utils.AnyOf(versions, issue.FixVersions...)
	}
}
