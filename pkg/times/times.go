package times

import (
	"time"
)

func ToJPTime(t time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	return t.In(loc)
}

func FormatToString(t time.Time) string {
	// e.g. 2025-02-10T12:04:26+09:00
	return t.Format(time.RFC3339)
}
