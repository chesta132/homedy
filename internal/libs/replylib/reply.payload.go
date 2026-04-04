package replylib

import (
	"time"

	"github.com/chesta132/goreply/reply"
)

// type for swaggo
type CodeError string

const (
	SwagCodeNotFound            CodeError = "NOT_FOUND"
	SwagCodeServerError         CodeError = "SERVER_ERROR"
	SwagCodeBadRequest          CodeError = "BAD_REQUEST"
	SwagCodeBadGateway          CodeError = "BAD_GATEWAY"
	SwagCodeUnauthorized        CodeError = "UNAUTHORIZED"
	SwagCodeConflict            CodeError = "CONFLICT"
	SwagCodeForbidden           CodeError = "FORBIDDEN"
	SwagCodeUnprocessableEntity CodeError = "UNPROCESSABLE_ENTITY"
	SwagCodeTooManyRequests     CodeError = "TOO_MANY_REQUESTS"
	SwagCodeServiceUnavailable  CodeError = "SERVICE_UNAVAILABLE"
	SwagCodeGatewayTimeout      CodeError = "GATEWAY_TIMEOUT"
	SwagCodeMethodNotAllowed    CodeError = "METHOD_NOT_ALLOWED"
	SwagCodeNotAcceptable       CodeError = "NOT_ACCEPTABLE"
	SwagCodeRequestTimeout      CodeError = "REQUEST_TIMEOUT"
	SwagCodePayloadTooLarge     CodeError = "PAYLOAD_TOO_LARGE"
	SwagCodeUnsupportedMedia    CodeError = "UNSUPPORTED_MEDIA_TYPE"
	SwagCodeGone                CodeError = "GONE"
	SwagCodeNotImplemented      CodeError = "NOT_IMPLEMENTED"
)

type status string

const (
	Success status = "SUCCESS"
	Error   status = "ERROR"
)

type Pagination struct {
	Current int  `json:"current"` // current offset
	HasNext bool `json:"hasHext"` // true if data more than replied
	Next    int  `json:"next"`    // if hasHext is false, next is 0
}

type Meta struct {
	Status      status            `json:"status" example:"SUCCESS"` // SUCCESS or ERROR
	Timestamp   time.Time         `json:"timestamp" example:"2006-01-02T15:04:05Z07:00"`
	Pagination  *Pagination       `json:"pagination,omitempty"`
	Information string            `json:"information,omitempty"`
	Tokens      map[string]string `json:"tokens,omitempty"`
	Debug       any               `json:"debug,omitempty"`
}

type Envelope struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`
}

func transformer(rp *reply.Reply) any {
	meta := rp.Meta()
	data := rp.Data()

	transMeta := Meta{
		Status:      status(meta.Status),
		Timestamp:   time.Unix(meta.Timestamp, 0),
		Information: meta.Info,
		Tokens:      meta.Tokens,
		Debug:       meta.Debug,
	}

	if meta.Pagination != nil {
		transMeta.Pagination = &Pagination{
			Current: meta.Pagination.Current,
			HasNext: meta.Pagination.HasNext,
			Next:    meta.Pagination.Next,
		}
	}

	return &Envelope{Data: data, Meta: transMeta}
}
