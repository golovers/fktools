package main

import (
	"flag"

	"github.com/golovers/fktools/fk"
)

func main() {
	plugin := flag.String("p", "jira", "one of the supported plugin: jira, phbricator, versionone")
	export := flag.String("e", "", "csv file name to be exported base on the base query")
	flag.Parse()

	var plug fk.Plugin
	switch *plugin {
	case "jira":
		plug = jira()
	default:
		panic("error: not supported plugin type")
	}
	fk.SetPlugin(plug)

	if *export != "" {
		issues, err := fk.AllIssues()
		if err != nil {
			panic(err)
		}
		issues.ToCSV(*export)
	}
}

func jira() fk.Plugin {
	var conf fk.JiraConf
	fk.LoadEnvConf(&conf)
	return fk.NewJira(&conf)
}
