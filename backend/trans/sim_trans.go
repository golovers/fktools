package trans

import (
	"github.com/golovers/kiki/backend/types"
	"github.com/golovers/kiki/backend/utils"
)

type simTransformer struct{}

func NewSimTrans() Transformer {
	return &simTransformer{}
}

func (tx *simTransformer) status(s string) string {
	if utils.Contains([]string{"Critical", "High", "Highest", "Blocker"}, s) {
		return "Critical"
	}
	if utils.Contains([]string{"Major"}, s) {
		return "Major"
	}
	if utils.Contains([]string{"Low", "Minor", "Medium", ""}, s) {
		return "Minor"
	}
	return s
}

func (tx *simTransformer) Transform(issue types.Issue) *types.Issue {
	newIssue := &issue
	newIssue.Status = tx.status(newIssue.Status)
	if newIssue.Priority == "" {
		newIssue.Priority = " "
	}
	return newIssue
}
