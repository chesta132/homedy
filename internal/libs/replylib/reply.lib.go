package replylib

import (
	"errors"

	"github.com/chesta132/goreply/reply"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func ErrorPayloadToErrorArg(err reply.ErrorPayload) (string, string, reply.ErrorOption, reply.ErrorOption) {
	return err.Code, err.Message, reply.WithDetails(err.Details), reply.WithFields(err.Fields)
}

func HandleError(err error, rp *reply.Reply) {
	if err, ok := err.(*reply.ErrorPayload); ok {
		rp.Error(ErrorPayloadToErrorArg(*err)).FailJSON()
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, redis.Nil) {
		rp.Error(CodeNotFound, err.Error()).FailJSON()
		return
	}
	rp.Error(CodeServerError, err.Error()).FailJSON()
}
