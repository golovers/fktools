package plug

import (
	"github.com/golovers/kiki/backend/types"
)

var plug Plugin

type Plugin interface {
	AllIssues() (types.Issues, error)
	CurrSprint() (string, error)
}

// SetPlugin set the plugin to JIRA, VersionOne, Phabricator,...
func SetPlugin(p Plugin) {
	plug = p
}

func AllIssues() (types.Issues, error) {
	return plug.AllIssues()
}

// CurrSprint return current sprint string like 'Sprint 18'
func CurrSprint() (string, error) {
	return plug.CurrSprint()
}
