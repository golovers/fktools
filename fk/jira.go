package fk

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

// JiraConf configurations for JIRA
type JiraConf struct {
	Host           string `envconfig:"JIRA_HOST"`
	Username       string `envconfig:"JIRA_USERNAME"`
	Password       string `envconfig:"JIRA_PASSWORD"`
	StoryPointAttr string `envconfig:"JIRA_ATTR_STORYPOINT"`
	SprintAttr     string `envconfig:"JIRA_ATTR_SPRINT"`
	EpicNameAttr   string `envconfig:"JIRA_ATTR_EPICLINK"`
	BaseQuery      string `envconfig:"JIRA_BASE_QUERY"`
}

// Jira plugin for JIRA
type Jira struct {
	conf *JiraConf
}

// NewJira return a new jira config
func NewJira(conf *JiraConf) *Jira {
	return &Jira{
		conf: conf,
	}
}

// AllIssues return all issues base on the configured BaseQuery
func (jr *Jira) AllIssues() (Issues, error) {
	issues := make([]*Issue, 0)
	start := 0
	max := 200
	for {
		jrIssues, res, err := jr.client().Issue.Search(jr.conf.BaseQuery, &jira.SearchOptions{
			StartAt:    start,
			MaxResults: max,
		})
		if err != nil {
			fmt.Println(res)
			return issues, err
		}
		for _, issue := range jrIssues {
			issues = append(issues, jr.toIssue(issue))
		}
		if len(jrIssues) < max {
			return issues, nil
		}
		start += max
	}
}

func (jr *Jira) toIssue(issue jira.Issue) *Issue {
	epic, _ := issue.Fields.Unknowns.String(jr.conf.EpicNameAttr)
	sprint := jr.sprint(issue)
	storyPoint := jr.floatVal(issue, jr.conf.StoryPointAttr)

	val, _ := issue.Fields.Unknowns.Value("versions")
	versions := make([]string, 0)
	for _, vv := range val.([]interface{}) {
		versions = append(versions, vv.(map[string]interface{})["name"].(string))
	}
	return &Issue{
		IssueType:              issue.Fields.Type.Name,
		Key:                    issue.Key,
		Status:                 issue.Fields.Status.Name,
		Summary:                issue.Fields.Summary,
		Priority:               issue.Fields.Priority.Name,
		StoryPoints:            storyPoint,
		EpicLink:               epic,
		Sprint:                 sprint,
		FixVersions:            toFixVersions(issue.Fields.FixVersions),
		Reporter:               issue.Fields.Reporter.Name,
		AffectsVersions:        versions,
		Assignee:               toAssignee(issue.Fields.Assignee),
		Components:             toComponents(issue.Fields.Components),
		Created:                StringToTime(issue.Fields.Created),
		Labels:                 issue.Fields.Labels,
		Resolution:             toResolution(issue.Fields.Resolution),
		Resolved:               StringToTime(issue.Fields.Resolutiondate),
		TotalOriginalEstimate:  timeOriginal(issue.Fields.TimeTracking),
		TotalRemainingEstimate: timeRemaining(issue.Fields.TimeTracking),
		TotalTimeSpent:         timeSpent(issue.Fields.TimeTracking),
	}
}

func (jr *Jira) client() *jira.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	jiraClient, err := jira.NewClient(client, jr.conf.Host)
	if err != nil {
		panic(err)
	}
	jiraClient.Authentication.SetBasicAuth(jr.conf.Username, jr.conf.Password)
	return jiraClient
}

func toComponents(components []*jira.Component) []string {
	cmps := make([]string, 0)
	for _, c := range components {
		cmps = append(cmps, c.Name)
	}
	return cmps
}

func toFixVersions(v []*jira.FixVersion) []string {
	versions := make([]string, 0)
	for _, ver := range v {
		versions = append(versions, ver.Name)
	}
	return versions
}

func (jr *Jira) floatVal(issue jira.Issue, attr string) float64 {
	v, ok := issue.Fields.Unknowns.Value(attr)
	if !ok {
		v = 0
	}
	rv := 0.0
	if v != nil {
		switch v.(type) {
		case int:
			rv = float64(reflect.ValueOf(v).Int())
		case float64:
			rv = reflect.ValueOf(v).Float()
		}
	}
	return rv
}

func toAssignee(assignee *jira.User) string {
	if assignee == nil {
		return ""
	}
	return assignee.Name
}

func toResolution(r *jira.Resolution) string {
	if r == nil {
		return ""
	}
	return r.Name
}

func timeSpent(t *jira.TimeTracking) int {
	if t == nil {
		return 0
	}
	return t.TimeSpentSeconds
}

func timeOriginal(t *jira.TimeTracking) int {
	if t == nil {
		return 0
	}
	return t.OriginalEstimateSeconds
}

func timeRemaining(t *jira.TimeTracking) int {
	if t == nil {
		return 0
	}
	return t.RemainingEstimateSeconds
}

func (jr *Jira) sprint(issue jira.Issue) string {
	vsprint, _ := issue.Fields.Unknowns.Value(jr.conf.SprintAttr)
	if vsprint != nil {
		val, _ := reflect.ValueOf(vsprint).Index(0).Interface().(string)
		vals := strings.Split(val, ",")
		for _, v := range vals {
			kv := strings.Split(v, "=")
			if kv[0] == "name" {
				return kv[1]
			}
		}
	}
	return ""
}
