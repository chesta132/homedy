package validatorlib

import (
	"reflect"
	"regexp"
)

var (
	usernameRegex = regexp.MustCompile(`^(?![0-9]+$)[a-zA-Z0-9_]{1,30}$|^([a-zA-Z0-9_](?:(?:[a-zA-Z0-9_]|(?:\.(?! \.))) {0,28}(?:[a-zA-Z0-9_]))?)$`)
	passwordRegex = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,32}$`)
)

func ValidateUsername(value any) bool {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t.Kind() != reflect.String {
		return false
	}
	return usernameRegex.MatchString(v.String())
}

func ValidatePassword(value any) bool {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t.Kind() != reflect.String {
		return false
	}
	return passwordRegex.MatchString(v.String())
}
