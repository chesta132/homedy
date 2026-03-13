package replylib

import "github.com/chesta132/goreply/reply"

func ErrorPayloadToErrorArg(err reply.ErrorPayload) (string, string, reply.ErrorOption, reply.ErrorOption) {
	return err.Code, err.Message, reply.WithDetails(err.Details), reply.WithFields(err.Fields)
}
