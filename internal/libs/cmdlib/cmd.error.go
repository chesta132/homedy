package cmdlib

import "errors"

var (
	ErrNoPkgToInstall     = errors.New("no pkg to install")
	ErrNoServiceToRestart = errors.New("no service to restart")
)
