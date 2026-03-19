package converter

import (
	"homedy/internal/libs/cmdlib"
	"homedy/internal/libs/logger"
)

func init() {
	_, err := cmdlib.InstallPkgs("libreoffice")
	// non no pkg to install fatal error
	if err != nil && err != cmdlib.ErrNoPkgToInstall {
		logger.Converter.Fatal(err.Error())
	}
}
