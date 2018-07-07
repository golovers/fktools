package main

import (
	"net/http"

	"github.com/golovers/fktools/fk/api"
	"github.com/golovers/fktools/fk/conf"
	"github.com/golovers/fktools/fk/db"
	"github.com/golovers/fktools/fk/iss"
	"github.com/golovers/fktools/fk/plug"
	"github.com/golovers/fktools/fk/rules"
	"github.com/golovers/fktools/fk/sched"
	"github.com/golovers/fktools/fk/trans"
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
	iss.SetIssueSvc(iss.NewSimIssueSvc())
	rules.SetRuleSvc(rules.NewSimRuleSvc())

	sched.SetScheduler(sched.NewCronScheduler())
	sched.Start()
	defer sched.Stop()

	api.SchedSync(conf.SyncSched)
	go iss.Sync()

	api.SchedRules()

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
