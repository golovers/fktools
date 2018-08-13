package issues

import (
	"encoding/json"

	"github.com/golovers/kiki/backend/db"
	"github.com/golovers/kiki/backend/plug"
	"github.com/golovers/kiki/backend/trans"
	"github.com/golovers/kiki/backend/types"
	"github.com/sirupsen/logrus"
)

var keyIssue = []byte("issues")

type simSvc struct {
}

func NewSimIssueSvc() IssueSvc {
	return &simSvc{}
}

// AllIssues return all iss
func (s *simSvc) Load() (types.Issues, error) {
	data, err := db.Get(keyIssue)
	if err != nil {
		return types.Issues{}, err
	}
	var issues types.Issues
	err = json.Unmarshal(data, &issues)
	if err != nil {
		return types.Issues{}, err
	}
	issues.Sort()
	return transform(issues), nil
}

// Sync sync data from remote to local for later use
func (s *simSvc) Sync() {
	logrus.Info("syncing data from remote...")
	issues, err := plug.AllIssues()
	if err != nil {
		logrus.Error("failed to sync data from remote", err)
		return
	}
	saveIssues(issues)
	logrus.Info("finished sync data from remote")
}

func saveIssues(issues types.Issues) {
	data, err := json.Marshal(issues)
	if err != nil {
		logrus.Errorf("failed to save iss: %v", err)
	}
	db.Put(keyIssue, data)
}

func transform(issues types.Issues) types.Issues {
	rs := make(types.Issues, 0)
	for _, issue := range issues {
		nw := trans.Transform(*issue)
		rs = append(rs, nw)
	}
	return rs
}
