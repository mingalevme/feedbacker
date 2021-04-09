package util

import (
	"github.com/mingalevme/feedbacker/pkg/strutils"
	"strings"
)

func Sprintf(template string, params map[string]interface{}) string {
	return strutils.Sprintf(template, params)
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

func IsNonEmptyString(s interface{}) bool {
	return strutils.IsNonEmptyString(s)
}

func IsEmptyString(s interface{}) bool {
	return strutils.IsEmptyString(s)
}

func IsPointerToEmptyString(s *string) bool {
	return strutils.IsPointerToEmptyString(s)
}

func IsPointerToNonEmptyString(s *string) bool {
	return strutils.IsPointerToNonEmptyString(s)
}