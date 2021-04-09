package strutils

import (
	"fmt"
	"strings"
)

func Sprintf(template string, params map[string]interface{}) string {
	for key, val := range params {
		template = strings.Replace(template, "%{"+key+"}s", fmt.Sprintf("%s", val), -1)
	}
	return template
}

func StrToPointerStr(str string) *string {
	return &str
}

func IsNonEmptyString(s interface{}) bool {
	return !IsEmptyString(s)
}

func IsEmptyString(s interface{}) bool {
	switch v := s.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case *string:
		return v == nil || strings.TrimSpace(*v) == ""
	default:
		panic("s must be string or *string")
	}
}

func IsPointerToEmptyString(s *string) bool {
	if s == nil {
		return true
	}
	if strings.TrimSpace(*s) == "" {
		return true
	}
	return false
}

func IsPointerToNonEmptyString(s *string) bool {
	return !IsPointerToEmptyString(s)
}
