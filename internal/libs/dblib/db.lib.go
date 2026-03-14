package dblib

import (
	"errors"
	"fmt"
	"homedy/internal/libs/replylib"
	"sync"

	"github.com/chesta132/goreply/reply"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GormErrorToReplyError(err error, dst any) *reply.ErrorPayload {
	s, err := schema.Parse(&dst, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		s.Table = "resource"
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &reply.ErrorPayload{
			Code:    replylib.CodeNotFound,
			Message: fmt.Sprintf("%s not found", s.Table),
		}
	}
	return &reply.ErrorPayload{
		Code:    replylib.CodeServerError,
		Message: err.Error(),
	}
}
