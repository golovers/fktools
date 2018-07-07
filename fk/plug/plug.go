package plug

import (
	"github.com/golovers/fktools/fk/types"
)

var plug Plugin

type Plugin interface {
	AllIssues() (types.Issues, error)
}

// SetPlugin set the plugin to JIRA, VersionOne, Phabricator,...
func SetPlugin(p Plugin) {
	plug = p
}

func AllIssues() (types.Issues, error) {
	return plug.AllIssues()
}
