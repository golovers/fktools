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
