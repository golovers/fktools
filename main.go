package main

import (
	"net/http"

	"github.com/golovers/kiki/api"
	"github.com/golovers/kiki/backend/conf"
	"github.com/golovers/kiki/backend/db"
	"github.com/golovers/kiki/backend/issues"
	"github.com/golovers/kiki/backend/plug"
	"github.com/golovers/kiki/backend/sched"
	"github.com/golovers/kiki/backend/trans"
)

func main() {
	conf, p := plugin()
	plug.SetPlugin(p)

	ldb, err := db.NewLDBDatabase(conf.DBName, 1000, 16)
	defer db.Close()
	if err != nil {
		panic(err)
	}
	db.SetDatabase(ldb)

	trans.SetTransformer(trans.NewSimTrans())
	issues.SetIssueSvc(issues.NewSimIssueSvc())

	sched.SetScheduler(sched.NewCronScheduler())
	sched.Start()
	defer sched.Stop()

	//api.SchedSync(conf.SyncSched)
	go issues.Sync()

	router := api.NewRouter()
	if err := http.ListenAndServe(conf.HTTPAddress, router); err != nil {
		panic(err)
	}
}

func plugin() (*conf.Conf, plug.Plugin) {
	var cfg conf.Conf
	conf.LoadEnvConf(&cfg)
	switch cfg.Plug {
	case "jira":
		return &cfg, plug.NewJira(&cfg)
	default:
		return &cfg, plug.NewJira(&cfg)
	}
}
