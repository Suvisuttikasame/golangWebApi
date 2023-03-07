package util

import "time"

func MockTime() time.Time {
	return time.Date(2023, 1, 1, 12, 45, 40, 3, time.UTC)
}
