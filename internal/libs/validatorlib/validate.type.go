package validatorlib

import "github.com/go-playground/validator/v10"

type Validatable interface {
	ValidateStruct(sl validator.StructLevel)
}
