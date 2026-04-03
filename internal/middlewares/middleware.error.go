package middlewares

import "errors"

var (
	ErrMiddlewareSkipped = errors.New("middeware skipped")
)
