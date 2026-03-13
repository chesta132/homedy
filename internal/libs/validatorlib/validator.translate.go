package validatorlib

import (
	"fmt"
	"homedy/internal/services/samba"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type translator func(fieldName string, err validator.FieldError) string

var translateErrorMap = map[string]translator{
	"email":            email,
	"oneof":            enum,
	"min":              length,
	"max":              length,
	"required":         required,
	"required_if":      requiredIf,
	"required_without": requiredWithout,
	"uuid4":            uuid4,
	"samba_bool":       createEnum(samba.Bools),
	"share_name":       shareName,
	"abs_path":         absolutePath,
	"file_permission":  filePerm,
}

func email(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%v is not valid email", err.Value())
}

func enum(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%v is not a valid enum of %s", err.Value(), err.Param())
}

func length(fieldName string, err validator.FieldError) string {
	suffix := ""
	switch err.Kind() {
	case reflect.Slice, reflect.Array:
		suffix = " items"
	case reflect.String:
		suffix = " characters"
	}

	operator := "at least"
	if err.Tag() == "max" {
		operator = "at most"
	}
	return fmt.Sprintf("%s must be %s %s%s", fieldName, operator, err.Param(), suffix)
}

func required(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%s is required", fieldName)
}

func requiredIf(fieldName string, err validator.FieldError) string {
	targetField, targetValue, _ := strings.Cut(err.Param(), " ")
	return fmt.Sprintf("%s is required if value of %s is %s", fieldName, targetField, targetValue)
}

func requiredWithout(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%s required if %s empty", fieldName, err.Param())
}

func uuid4(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%s must be an uuid4", fieldName)
}

func createEnum[E ~string](enum []E) translator {
	return func(fieldName string, err validator.FieldError) string {
		return fmt.Sprintf("%v is not a valid enum of %s", err.Value(), enum)
	}
}

// samba

func shareName(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%v is not a valid share name", err.Value())
}

// iolib

func absolutePath(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%v is not a valid absolute path", err.Value())
}

func filePerm(fieldName string, err validator.FieldError) string {
	return fmt.Sprintf("%v is not a valid file permission", err.Value())
}
