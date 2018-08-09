package reports

import (
	"strings"

	"github.com/golovers/kiki/backend/types"
	"github.com/golovers/kiki/backend/utils"
)

type FilterFunc func(issue *types.Issue) bool

var defectFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "defect", "bug")
}

var storyFilter = func(issue *types.Issue) bool {
	return utils.OneOf(strings.ToLower(issue.IssueType), "story", "improvement", "enhancement")
}

var teamSprintFilter = func(issueTypeFilter FilterFunc, sprint string, team string) FilterFunc {
	return func(issue *types.Issue) bool {
		if !issueTypeFilter(issue) {
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

var teamSprintStoryFilter = func(sprint string, team string) FilterFunc {
	return teamSprintFilter(storyFilter, sprint, team)
}

var teamSprintDefectFilter = func(sprint string, team string) FilterFunc {
	return teamSprintFilter(defectFilter, sprint, team)
}

var otherTeamsSprintFilter = func(issueTypeFilter FilterFunc, sprint string, teams []string) FilterFunc {
	return func(issue *types.Issue) bool {
		if !issueTypeFilter(issue) {
			return false
		}
		for _, team := range teams {
			if utils.OneOf(team, issue.Labels...) {
				return false
			}
		}
		if sprint != "*" && strings.ToLower(sprint) != strings.ToLower(issue.Sprint) {
			return false
		}
		return true
	}
}

var otherTeamsSprintDefectFilter = func(sprint string, teams []string) FilterFunc {
	return otherTeamsSprintFilter(defectFilter, sprint, teams)
}

var otherTeamsSprintStoryFilter = func(sprint string, teams []string) FilterFunc {
	return otherTeamsSprintFilter(storyFilter, sprint, teams)
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
