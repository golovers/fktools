package trans

import (
	"github.com/golovers/kiki/backend/types"
)

var tx Transformer

type Transformer interface {
	Transform(issue types.Issue) *types.Issue
}

func SetTransformer(t Transformer) {
	tx = t
}

func Transform(issue types.Issue) *types.Issue {
	return tx.Transform(issue)
}
