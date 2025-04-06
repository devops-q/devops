package utils

import (
	"time"
)

func FormatDateTime(t time.Time) string {
	return t.UTC().Format("2006-01-02 @ 15:04")
}

func ToISODateString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
