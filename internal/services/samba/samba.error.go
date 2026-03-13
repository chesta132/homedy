package samba

import "errors"

var (
	ErrShareAlreadyExist   = errors.New("share config with this name already exist")
	ErrShareNotExist       = errors.New("share config with this name does not exist")
	ErrPathAlreadyExist    = errors.New("folder sharing with this path already exist")
	ErrConfigNotExist      = errors.New("share configuration does not exist")
	ErrNameContainsInvalid = errors.New("share name contains invalid characters/name")
)
