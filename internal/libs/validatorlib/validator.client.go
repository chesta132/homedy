package validatorlib

import (
	"homedy/internal/services/samba"
	"reflect"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Client = validator.New(validator.WithRequiredStructEnabled())

func registFirstTag(tags map[string]string) validator.TagNameFunc {
	return func(fld reflect.StructField) string {
		for tag, splitter := range tags {
			name, _, _ := strings.Cut(fld.Tag.Get(tag), splitter)
			if name != "-" && name != "" {
				return name
			}
		}
		return ""
	}
}

func registEnumValidation[E ~string](enum []E) validator.Func {
	return func(fl validator.FieldLevel) bool {
		if fl.Field().Kind() != reflect.String {
			return false
		}

		value := fl.Field().String()
		if !slices.Contains(enum, E(value)) {
			return false
		}

		return true
	}
}

func basicValidatorToValidatorFunc(f func(any) bool) validator.Func {
	return func(fl validator.FieldLevel) bool {
		return f(fl.Field().Interface())
	}
}

func init() {
	// register tags
	Client.RegisterTagNameFunc(registFirstTag(map[string]string{
		"json": ",",
		"uri":  ",",
		"form": ",",
	}))

	// register enum
	Client.RegisterValidation("samba_bool", registEnumValidation(samba.Bools))

	// register basic validator
	Client.RegisterValidation("share_name", basicValidatorToValidatorFunc(ValidateShareName))
	Client.RegisterValidation("abs_path", basicValidatorToValidatorFunc(ValidateAbsPath))
	Client.RegisterValidation("file_permission", basicValidatorToValidatorFunc(ValidateFilePermission))
}
