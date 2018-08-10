package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golovers/kiki/backend/reports"
)

// Index write index info
func Index(w http.ResponseWriter, r *http.Request) {
	sprint, teams := sprintAndTeamFromRequest(r)
	if sprint == "" {
		sprint = "*"
	}
	status := reports.Sprint(sprint, teams...)
	if sprint == "*" {
		status.Sprint = "All Sprints"
	}
	indexTmpl.Execute(w, r, status)
}

// DefectStatus provide defect status over entire the backlog
func DefectStatus(w http.ResponseWriter, r *http.Request) {
	status := reports.Defects()
	b, err := json.Marshal(status)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to marshal")
	}
	w.Write(b)
}

//StoryStatus provide stories status over entire the backlog
func StoryStatus(w http.ResponseWriter, r *http.Request) {
	status := reports.Stories()
	b, err := json.Marshal(status)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to marshal")
		return
	}
	w.Write(b)
}

//SprintStatus provide status of sprints which will have details report for each teams
func SprintStatus(w http.ResponseWriter, r *http.Request) {
	sprint, teams := sprintAndTeamFromRequest(r)
	if sprint == "" {
		writeErr(w, http.StatusInternalServerError, "sprint is required")
		return
	}
	status := reports.Sprint(sprint, teams...)
	data, err := json.Marshal(status)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "failed to marshal")
		return
	}
	w.Write(data)
}

func sprintAndTeamFromRequest(r *http.Request) (sprint string, teams []string) {
	vars := r.URL.Query()
	sprint = vars.Get("sprint")
	teams = strings.Split(vars.Get("teams"), ",")
	return
}

func writeErr(w http.ResponseWriter, code int, err string) {
	w.WriteHeader(code)
	er := errMsg{
		Code:  code,
		Error: err,
	}
	data, _ := json.Marshal(er)
	w.Write(data)
}

type errMsg struct {
	Code  int
	Error string
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}
