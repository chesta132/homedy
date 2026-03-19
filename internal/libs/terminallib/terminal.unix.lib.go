//go:build linux || darwin

package terminallib

import (
	"os/user"
	"strconv"
	"syscall"
)

func GetUserCredential(u *user.User) (*syscall.SysProcAttr, error) {
	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return nil, err
	}

	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return nil, err
	}

	return &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}, nil
}
