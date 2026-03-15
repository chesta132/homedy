package sambalib

import (
	"homedy/internal/libs/cmdlib"
	"homedy/internal/libs/logger"
)

func init() {
	cmd, err := cmdlib.InstallPkgs("samba")
	// non no pkg to install fatal error
	if err != nil && err != cmdlib.ErrNoPkgToInstall {
		logger.Samba.Fatal(err.Error())
	}
	// backup conf
	if cmd != nil && err == nil {
		err = Backup()
		if err != nil {
			logger.Samba.Fatal(err.Error(), logger.Fields("rec_step", "manual copy your smb conf"))
		}
	}
}
