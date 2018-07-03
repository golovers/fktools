package fk

import (
	"io/ioutil"
	"time"

	"strings"

	"github.com/go-yaml/yaml"
	"github.com/kelseyhightower/envconfig"
)

// ReadConf read config from  file and return config struct
func ReadConf(file string, cfg interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err1 := yaml.Unmarshal([]byte(data), &cfg)
	if err1 != nil {
		panic(err1)
	}
}

// StringToTime convert string to time
// TODO implement me
func StringToTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05.000+0200", strings.Trim(s, " "))
	if err != nil {
		panic(err)
	}
	return &t
}

func contains(values []string, v string) bool {
	for _, val := range values {
		if val == v {
			return true
		}
	}
	return false
}

// LoadEnvConf load configurations from config file
func LoadEnvConf(t interface{}) {
	if err := envconfig.Process("", t); err != nil {
		panic(err)
	}
}
