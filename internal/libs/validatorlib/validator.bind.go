package validatorlib

import "github.com/gin-gonic/gin"

type binder func(any) error

func BindAndValidate[T any](binders ...binder) (T, error) {
	var toBind T
	for _, binder := range binders {
		if err := binder(&toBind); err != nil {
			return toBind, err
		}
	}
	if errPayload := ValidateStructToReply(toBind); errPayload != nil {
		return toBind, errPayload
	}
	return toBind, nil
}

func BindJSONAndValidate[T any](c *gin.Context) (T, error) {
	return BindAndValidate[T](c.ShouldBindJSON)
}
