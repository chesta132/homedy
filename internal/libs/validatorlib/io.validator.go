package validatorlib

import (
	"reflect"
	"regexp"
)

var absPathRegex = regexp.MustCompile("^(/[^/ ]*)+/?$")

func ValidateAbsPath(value any) bool {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t.Kind() != reflect.String {
		return false
	}
	return absPathRegex.MatchString(v.String())
}

func ValidateFilePermission(value any) bool {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false
		}
		return ValidateFilePermission(v.Elem().Interface())
	}

	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		return false
	}

	vLen := v.Len()
	if vLen != 3 {
		return false
	}

	for i := 0; i < vLen; i++ {
		val := v.Index(i)
		if val.Kind() != reflect.Int {
			return false
		}
		valInt := val.Int()
		if valInt < 0 || valInt > 7 {
			return false
		}
	}

	return true
}
