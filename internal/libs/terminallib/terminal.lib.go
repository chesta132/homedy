package terminallib

import (
	"os/exec"
	"os/user"
)

func EnsureUser(username string) error {
	_, err := user.Lookup(username)
	if err == nil {
		return nil
	}

	cmd := exec.Command("useradd", "-m", username)
	return cmd.Run()
}
