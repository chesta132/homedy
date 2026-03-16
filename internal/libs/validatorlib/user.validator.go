package validatorlib

import (
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]([a-zA-Z0-9_.]{1,28}[a-zA-Z0-9_]|[a-zA-Z0-9_]?)$`)
	passwordRegex = regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,32}$`)
	allDigits     = regexp.MustCompile(`^[0-9]+$`)
)

func ValidateUsername(value any) bool {
	str, ok := validateStr(value)
	if !ok {
		return false
	}

	if !usernameRegex.MatchString(str) {
		return false
	}
	return !allDigits.MatchString(str)
}

func ValidatePassword(value any) bool {
	str, ok := validateStr(value)
	if !ok {
		return false
	}

	if !passwordRegex.MatchString(str) {
		return false
	}

	// go doesnt support this action
	var hasLower, hasUpper, hasDigit bool
	for _, c := range str {
		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}
	return hasLower && hasUpper && hasDigit
}
