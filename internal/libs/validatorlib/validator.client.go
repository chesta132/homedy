package validatorlib

import (
	"homedy/internal/models"
	"homedy/internal/models/payloads"

	"github.com/go-playground/validator/v10"
)

var Client = validator.New(validator.WithRequiredStructEnabled())

func init() {
	// register tags
	Client.RegisterTagNameFunc(registFirstTag(map[string]string{
		"json": ",",
		"uri":  ",",
		"form": ",",
	}))

	// register enum
	Client.RegisterValidation("samba_bool", registEnumValidation(models.SambaBools))

	// register basic validator
	Client.RegisterValidation("share_name", basicValidatorToValidatorFunc(ValidateShareName))
	Client.RegisterValidation("abs_path", basicValidatorToValidatorFunc(ValidateAbsPath))
	Client.RegisterValidation("file_permission", basicValidatorToValidatorFunc(ValidateFilePermission))
	Client.RegisterValidation("username", basicValidatorToValidatorFunc(ValidateUsername))
	Client.RegisterValidation("password", basicValidatorToValidatorFunc(ValidatePassword))

	registerValidatable(&payloads.RequestConvertMultiple{})
}
