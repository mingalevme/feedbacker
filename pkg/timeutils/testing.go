// +build testing

package timeutils

import "time"

func SetTestNow(now time.Time) {
	testNow = &now
}

func ResetTestNow() {
	testNow = nil
}
