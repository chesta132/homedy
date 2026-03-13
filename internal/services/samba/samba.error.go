package samba

import "errors"

var (
	ErrShareAlreadyExist = errors.New("share config with this name already exist")
	ErrShareNotExist     = errors.New("share config with this name is not exist")
	ErrPathAlreadyExist  = errors.New("folder sharing with this path already exist")
)
