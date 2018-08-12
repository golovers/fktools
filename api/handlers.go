package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/golovers/kiki/backend/links"
	"github.com/golovers/kiki/backend/reports"
)

func index(w http.ResponseWriter, r *http.Request) {
	ls := links.Links()
	logrus.Info("links: ", ls)
	indexTmpl.Execute(w, r, ls)
}

func linkConfig(w http.ResponseWriter, r *http.Request) {
	linksTmpl.Execute(w, r, links.Links())
}

func addLink(w http.ResponseWriter, r *http.Request) {
	link := &links.QuickLink{
		Name:     r.FormValue("name"),
		Link:     fmt.Sprintf("group?epic=%s&fixVersions=%s&labels=%s&sprint=%s", r.FormValue("epic"), r.FormValue("fixVersions"), r.FormValue("labels"), r.FormValue("sprint")),
		Visitted: 0,
	}
	logrus.Infof("Link - name: %s, link: %s", link.Name, link.Link)
	err := links.Add(link)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.Redirect(w, r, "/links", http.StatusFound)
}

func deleteLinks(w http.ResponseWriter, r *http.Request) {
	links.DeleteAll()
}

func sprintStatus(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	sprint := vars.Get("sprint")
	teams := strings.Split(vars.Get("teams"), ",")
	if sprint == "" {
		sprint = "*"
	}
	status := reports.Sprint(sprint, teams...)
	if sprint == "*" {
		status.Sprint = "All Sprints"
	}
	sprintTmpl.Execute(w, r, status)
}

func groupStatus(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	sprint := vars.Get("sprint")
	labels := strings.Split(vars.Get("labels"), ",")
	epic := vars.Get("epic")
	fixVersions := strings.Split(vars.Get("fixVersions"), ",")
	if sprint == "" {
		sprint = "*"
	}
	if len(labels) == 0 {
		labels = append(labels, "*")
	}
	if len(fixVersions) == 0 {
		fixVersions = append(fixVersions, "*")
	}
	if epic == "" {
		epic = "*"
	}
	logrus.Info(epic, fixVersions, labels, sprint)
	status := reports.EpicStatus(epic, fixVersions, labels, sprint)
	epicTmpl.Execute(w, r, status)
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
