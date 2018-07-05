package fk

import "github.com/kelseyhightower/envconfig"

// Conf configurations
type Conf struct {
	Plug string `envconfig:"FK_PLUGIN" default:"jira"`

	Host          string            `envconfig:"FK_HOST"`
	Username      string            `envconfig:"FK_USERNAME"`
	Password      string            `envconfig:"FK_PASSWORD"`
	FieldsMapping map[string]string `envconfig:"FK_FIELD_MAPPING"`
	BaseQuery     string            `envconfig:"FK_BASE_QUERY"`

	HTTPAddress string `envconfig:"FK_HTTP_ADDRESS" default:":8080"`
	SyncSched   string `envconfig:"FK_SYNC_SCHED" default:"@every 5m"`
	DBName      string `envconfig:"FK_DB_NAME" default:"fk.db"`
}

// LoadEnvConf load configurations from config file
func LoadEnvConf(t interface{}) {
	if err := envconfig.Process("", t); err != nil {
		panic(err)
	}
}
