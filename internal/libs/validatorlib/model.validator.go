package validatorlib

import (
	"reflect"
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]([a-zA-Z0-9_.]{1,28}[a-zA-Z0-9_]|[a-zA-Z0-9_]?)$`)
	passwordRegex = regexp.MustCompile(`^[A-Za-z\d@$!%*?&]{8,32}$`)
	allDigits = regexp.MustCompile(`^[0-9]+$`)
)

func ValidateUsername(value any) bool {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t.Kind() != reflect.String {
		return false
	}

	vStr := v.String()
	if !usernameRegex.MatchString(vStr) {
		return false
	}
	return !allDigits.MatchString(vStr)
}

func ValidatePassword(value any) bool {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t.Kind() != reflect.String {
		return false
	}

	valStr := v.String()
	if !passwordRegex.MatchString(valStr) {
		return false
	}

	// go doesnt support this action
	var hasLower, hasUpper, hasDigit bool
	for _, c := range valStr {
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
