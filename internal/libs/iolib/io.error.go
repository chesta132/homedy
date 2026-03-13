package iolib

import "errors"

var (
	ErrPermissionLength   = errors.New("invalid permission length")
	ErrPermissionNotKnown = errors.New("permission numeric not known")
)
