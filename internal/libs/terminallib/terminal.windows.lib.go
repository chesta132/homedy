//go:build windows

package terminallib

import (
	"os/user"
	"syscall"
)

func GetUserCredential(u *user.User) (*syscall.SysProcAttr, error) {
	return &syscall.SysProcAttr{}, nil
}
