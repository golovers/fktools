package conf

import (
	"github.com/kelseyhightower/envconfig"
)

// Conf configurations
type Conf struct {
	Plug string `envconfig:"KIKI_PLUGIN" default:"jira"`

	Host          string            `envconfig:"KIKI_HOST"`
	Username      string            `envconfig:"KIKI_USERNAME"`
	Password      string            `envconfig:"KIKI_PASSWORD"`
	FieldsMapping map[string]string `envconfig:"KIKI_FIELD_MAPPING" default:"sprint:customfield_10801,epic:customfield_11209,storypoint:customfield_10011"`
	BaseQuery     string            `envconfig:"KIKI_BASE_QUERY" default:"issuetype not in (Task, Sub-task, Test)"`

	HTTPAddress string `envconfig:"KIKI_HTTP_ADDRESS" default:":8080"`
	SyncSched   string `envconfig:"KIKI_SYNC_SCHED" default:"@every 5m"`
	DBName      string `envconfig:"KIKI_DB_NAME" default:"kiki.db"`
}

// LoadEnvConf load configurations from config file
func LoadEnvConf(t interface{}) {
	if err := envconfig.Process("", t); err != nil {
		panic(err)
	}
}
