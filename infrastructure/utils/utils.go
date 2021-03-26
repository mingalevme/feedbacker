package utils

import (
	"fmt"
	"strings"
)

// https://play.golang.org/p/kn24JK87VI
func Sprintf(template string, params map[string]interface{}) string {
	for key, val := range params {
		template = strings.Replace(template, "%{"+key+"}s", fmt.Sprintf("%s", val), -1)
	}
	return template
}
