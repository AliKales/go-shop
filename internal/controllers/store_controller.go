package controllers

import (
	"regexp"
)

func IsLinkNameOK(val string) bool {
	pattern := "^[a-zA-Z0-9 ]*$"
	_, err := regexp.MatchString(pattern, val)

	if err != nil {
		return false
	}

	if len(val) > 25 {
		return false
	}

	return true
}
