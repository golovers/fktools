package fk

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"strings"

	"github.com/go-yaml/yaml"
	"github.com/rs/xid"
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

// IssuesToJSONString convert list of issues to a json string
func IssuesToJSONString(issues Issues) (string, error) {
	data, err := json.Marshal(issues)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GenID return a new uqique ID
func GenID() string {
	return xid.New().String()
}
