package sambalib

import (
	"fmt"
	"homedy/internal/libs/cmdlib"
	"homedy/internal/models"
)

func FilterShares(shares models.Shares) models.Shares {
	result := make(models.Shares)
	for k, v := range shares {
		if k != "global" && k != "printers" && k != "print$" {
			result[k] = v
		}
	}

	return result
}

func IsPathExist(shares models.Shares, share models.Share) bool {
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
