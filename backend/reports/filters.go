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

var oneOfTheseTeams = func(teams ...string) FilterFunc {
	return func(issue *types.Issue) bool {
		return utils.AnyOf(teams, issue.Labels...)
	}
}

var notTheseTeamsFilter = func(teams ...string) FilterFunc {
	return reverseFilter(oneOfTheseTeams(teams...))
}

var sprintFilter = func(sprint string) FilterFunc {
	return func(issue *types.Issue) bool {
		return sprint == "*" || utils.OneOf(sprint, issue.Sprint)
	}
}

var teamSprintFilter = func(issueTypeFilter FilterFunc, sprint string, team string) FilterFunc {
	return multipleFilters(issueTypeFilter, sprintFilter(sprint), oneOfTheseTeams(team))
}

var teamSprintStoryFilter = func(sprint string, team string) FilterFunc {
	return teamSprintFilter(storyFilter, sprint, team)
}

var teamSprintDefectFilter = func(sprint string, team string) FilterFunc {
	return teamSprintFilter(defectFilter, sprint, team)
}

var otherTeamsSprintFilter = func(issueTypeFilter FilterFunc, sprint string, teams []string) FilterFunc {
	return multipleFilters(issueTypeFilter, sprintFilter(sprint), notTheseTeamsFilter(teams...))
}

var otherTeamsSprintDefectFilter = func(sprint string, teams []string) FilterFunc {
	return otherTeamsSprintFilter(defectFilter, sprint, teams)
}

var otherTeamsSprintStoryFilter = func(sprint string, teams []string) FilterFunc {
	return otherTeamsSprintFilter(storyFilter, sprint, teams)
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
