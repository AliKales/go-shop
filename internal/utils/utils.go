package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsExpired(date time.Time) bool {
	return time.Now().UTC().After(date)
}

func IsOnlyNumbers(input string) bool {
	pattern := "^[0-9]+$"

	regexpPattern := regexp.MustCompile(pattern)

	return regexpPattern.MatchString(input)
}

func IsStringInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func LowerFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func LowerMapKeys(m map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, val := range m {
		lowerKey := LowerFirstLetter(key)
		result[lowerKey] = val
	}

	return result
}