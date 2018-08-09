package trans

import (
	"strings"

	"github.com/golovers/kiki/backend/types"
	"github.com/golovers/kiki/backend/utils"
)

type simTransformer struct{}

func NewSimTrans() Transformer {
	return &simTransformer{}
}

func (tx *simTransformer) priority(s string) string {
	s = strings.ToLower(s)
	if utils.Contains([]string{"critical", "high", "highest", "blocker"}, s) {
		return "critical"
	}
	if utils.Contains([]string{"major"}, s) {
		return "major"
	}
	if utils.Contains([]string{"low", "minor", "medium", "", "unclassified"}, s) {
		return "minor"
	}
	return s
}

func (tx *simTransformer) status(s string) string {
	s = strings.ToLower(s)
	if utils.OneOf(s, "", "open", "submitted", "nil") {
		return "open"
	}
	if utils.OneOf(s, "inprogress", "in progress", "in-progress", "on-going", "reviewing", "code-review", "code-reviewing") {
		return "inprogress"
	}
	if utils.OneOf(s, "resolved", "code-completed", "finished", "complete") {
		return "resolved"
	}
	if utils.OneOf(s, "reopened", "failed") {
		return "reopened"
	}
	if utils.OneOf(s, "closed", "done", "rejected", "invalid") {
		return "closed"
	}
	return s
}

func (tx *simTransformer) Transform(issue types.Issue) *types.Issue {
	newIssue := &issue
	newIssue.Priority = tx.priority(newIssue.Priority)
	newIssue.Status = tx.status(newIssue.Status)
	return newIssue
}
