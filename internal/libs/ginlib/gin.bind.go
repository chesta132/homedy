package ginlib

import (
	"homedy/internal/libs/validatorlib"
	"reflect"

	"github.com/gin-gonic/gin"
)

type Binder func(any) error

func Bind[T any](binders ...Binder) (T, error) {
	var toBind T
	for _, binder := range binders {
		if err := binder(&toBind); err != nil {
			return toBind, err
		}
	}
	return toBind, nil
}

func BindAndValidate[T any](binders ...Binder) (T, error) {
	toBind, err := Bind[T](binders...)
	if err != nil {
		return toBind, err
	}

	t := reflect.TypeOf(toBind)
	if t.Kind() == reflect.Struct {
		if errPayload := validatorlib.ValidateStructToReply(toBind); errPayload != nil {
			return toBind, errPayload
		}
	}

	return toBind, nil
}

func BindJSONAndValidate[T any](c *gin.Context) (T, error) {
	return BindAndValidate[T](c.ShouldBindJSON)
}
