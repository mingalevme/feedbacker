package timeutils

import "time"

var testNow *time.Time

func Now() time.Time {
	if testNow != nil {
		return *testNow
	}
	return time.Now()
}
