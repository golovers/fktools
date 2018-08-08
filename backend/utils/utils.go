package utils

import (
	"time"

	"strings"

	"github.com/rs/xid"
)

// StringToTime convert string to time
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

func Contains(values []string, v string) bool {
	for _, val := range values {
		if val == v {
			return true
		}
	}
	return false
}

// GenID return a new uqique ID
func GenID() string {
	return xid.New().String()
}

func StringToStrings(s string) []string {
	return strings.Split(s, ",")
}

func StringsToString(s []string) string {
	return strings.Join(s, ",")
}

func EscapeSpecialChars(s string) string {
	return strings.Replace(s, "\"", "\\\"", -1)
}

func OneOf(v string, values ...string) bool {
	for _, val := range values {
		if v == val {
			return true
		}
	}
	return false
}

// MergeFloat64Maps merge multiple maps together
func MergeFloat64Maps(maps ...map[string]float64) map[string]float64 {
	rs := make(map[string]float64)
	for _, m := range maps {
		for k, v := range m {
			if _, ok := rs[k]; ok {
				rs[k] += v
				continue
			}
			rs[k] = v
		}
	}
	return rs
}
