package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/sirupsen/logrus"

	"github.com/golovers/kiki/backend/links"
	"github.com/golovers/kiki/backend/reports"
)

func index(w http.ResponseWriter, r *http.Request) {
	ls := links.Links()
	logrus.Info("links: ", ls)
	indexTmpl.Execute(w, r, ls)
}

func groupLinks(w http.ResponseWriter, r *http.Request) {
	linksTmpl.Execute(w, r, links.LinksByType("group"))
}

func addLink(w http.ResponseWriter, r *http.Request) {
	link := linkFromRequest(r)
	logrus.Infof("Link - name: %s, link: %s", link.Name, link.Link)
	err := links.Add(link)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if link.Type == "group" {
		http.Redirect(w, r, "/links", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/sprint_links", http.StatusFound)
}

func linkFromRequest(r *http.Request) *links.QuickLink {
	var link *links.QuickLink
	if r.FormValue("type") == "group" {
		link = &links.QuickLink{
			Name:     r.FormValue("name"),
			Type:     "group",
			Link:     fmt.Sprintf("group?epic=%s&fixVersions=%s&labels=%s&sprint=%s", r.FormValue("epic"), r.FormValue("fixVersions"), r.FormValue("labels"), r.FormValue("sprint")),
			Visitted: 0,
		}
	} else {
		link = &links.QuickLink{
			Name: r.FormValue("name"),
			Type: "sprint",
			Link: fmt.Sprintf("sprint?sprint=%s&teams=%s", r.FormValue("sprint"), r.FormValue("teams")),
		}
	}
	return link
}

func runLink(w http.ResponseWriter, r *http.Request) {
	link := linkFromRequest(r)
	http.Redirect(w, r, link.Link, http.StatusFound)
}

func sprintLinks(w http.ResponseWriter, r *http.Request) {
	sprintLInksTmpl.Execute(w, r, links.LinksByType("sprint"))
}

func deleteLink(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	logrus.Info("links to be deleted: ", id)
	links.Delete(id)
	if mux.Vars(r)["type"] == "group" {
		http.Redirect(w, r, "/links", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/sprint_links", http.StatusFound)
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
