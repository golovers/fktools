package fk

import (
	"fmt"
	"os"
	"strings"
)

var (
	plugin      Plugin
	transformer Transformer
)

// SetPlugin set the plugin to JIRA, VersionOne, Phabricator,...
func SetPlugin(p Plugin) {
	plugin = p
}

// SetTransformer set the StatusTransformer to be used
func SetTransformer(s Transformer) {
	transformer = s
}

// AllIssues return all issues
func AllIssues() (Issues, error) {
	issues, err := plugin.AllIssues()
	if err != nil {
		return Issues{}, err
	}
	transform(issues)
	return issues, nil
}

func transform(issues Issues) {
	if transformer == nil {
		transformer = &CommonTransformer{}
	}
	for _, issue := range issues {
		issue = transformer.Transform(*issue)
	}
}

// ToCSV export issues into csv file
func (issue *Issues) ToCSV(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s\n", strings.Join(attrs, ",")))
	issues, err := AllIssues()
	if err != nil {
		return err
	}
	for i, issue := range issues {
		f.WriteString(issue.ToCSV())
		if i%100 == 0 {
			f.Sync()
		}
	}
	return nil
}
