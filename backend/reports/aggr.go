package reports

import "github.com/golovers/kiki/backend/types"

//AggrFunc calculate aggregation value
type AggrFunc func(issue *types.Issue) float64

var aggrCount = func(issue *types.Issue) float64 {
	return 1.0
}
var aggrStoryPoints = func(issue *types.Issue) float64 {
	return issue.StoryPoints
}
