package api

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

var (
	indexTmpl       = parseTemplate("index.html")
	epicTmpl        = parseTemplate("group.html")
	sprintTmpl      = parseTemplate("sprint.html")
	linksTmpl       = parseTemplate("links.html")
	sprintLInksTmpl = parseTemplate("sprint_links.html")
)

// parseTemplate applies a given file to the body of the base template.
func parseTemplate(filename string) *appTemplate {
	tmpl := template.Must(template.ParseFiles("templates/base.html"))
	funcMap := template.FuncMap{
		"sum": sum,
	}

	// Put the named file into a template called "body"
	path := filepath.Join("templates", filename)
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("could not read template: %v", err))
	}
	template.Must(tmpl.New("body").Funcs(funcMap).Parse(string(b)))

	return &appTemplate{tmpl.Lookup("base.html")}
}

// appTemplate is a user login-aware wrapper for a html/template.
type appTemplate struct {
	t *template.Template
}

// Execute writes the template using the provided data, adding login and user
// information to the base template.
func (tmpl *appTemplate) Execute(w http.ResponseWriter, r *http.Request, data interface{}) *errMsg {
	d := struct {
		Data interface{}
	}{
		Data: data,
	}
	if err := tmpl.t.Execute(w, d); err != nil {
		return &errMsg{http.StatusInternalServerError, err.Error()}
	}
	return nil
}

func sum(vals ...float64) float64 {
	v := 0.0
	for _, val := range vals {
		v += val
	}
	return v
}
