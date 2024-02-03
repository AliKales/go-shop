package utils

import (
	"regexp"
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