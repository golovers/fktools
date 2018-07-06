package main

import (
	"net/http"

	"github.com/golovers/fktools/fk"
	"github.com/gorilla/mux"
)

func main() {
	conf, plug := plugin()
	fk.SetPlugin(plug)

	db, err := fk.NewLDBDatabase(conf.DBName, 1000, 16)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	fk.SetDatabase(db)

	sched := fk.NewCronScheduler()
	sched.Start()
	fk.SetScheduler(sched)
	defer sched.Stop()

	fk.SchedSync(conf.SyncSched)
	go fk.Sync()

	fk.SchedRules()

	r := mux.NewRouter()
	if err := http.ListenAndServe(conf.HTTPAddress, r); err != nil {
		panic(err)
	}
}

func plugin() (*fk.Conf, fk.Plugin) {
	var conf fk.Conf
	fk.LoadEnvConf(&conf)
	switch conf.Plug {
	case "jira":
		return &conf, fk.NewJira(&conf)
	default:
		return &conf, fk.NewJira(&conf)
	}
}
