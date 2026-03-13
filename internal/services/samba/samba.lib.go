package samba

import (
	"fmt"
	"homedy/internal/libs/cmdlib"
)

func FilterShares(shares Shares) Shares {
	result := make(Shares)
	for k, v := range shares {
		if k != "global" && k != "printers" && k != "print$" {
			result[k] = v
		}
	}

	return result
}

func isPathExist(shares Shares, share Share) bool {
	for _, _share := range shares {
		if _share.Path == share.Path {
			return true
		}
	}
	return false
}

func restartService() error {
	cmd, err := cmdlib.RestartService("smbd", "nmbd")
	if err != nil {
		out, _ := cmd.Output()
		return fmt.Errorf("%w: %s", err, out)
	}
	return err
}
