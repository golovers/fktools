package iss

import (
	"encoding/json"

	"github.com/golovers/fktools/fk/db"
	"github.com/golovers/fktools/fk/plug"
	"github.com/golovers/fktools/fk/trans"
	"github.com/golovers/fktools/fk/types"
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
	transform(issues)
	return issues, nil
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

func transform(issues types.Issues) {
	for _, issue := range issues {
		issue = trans.Transform(*issue)
	}
}
