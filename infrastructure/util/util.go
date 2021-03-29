package util

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

func ParseHeader(header string) map[string]*string {
	m := make(map[string]*string)
	for _, pair := range strings.Split(header, ";") {
		values := strings.SplitN(pair, "=", 2)
		k := strings.TrimSpace(values[0])
		if len(values) == 1 {
			m[k] = nil
		} else {
			value := strings.TrimSpace(values[1])
			m[k] = &value
		}
	}
	return m
}
