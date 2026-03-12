package cmdlib

import (
	"homedy/internal/libs/logger"
	"os/exec"
	"strings"
	"time"
)

func InstallPkgs(pkgs ...string) (*exec.Cmd, error) {
	pkgsToInstall := filterUninstalledPkgs(pkgs)
	if len(pkgsToInstall) <= 0 {
		return nil, ErrNoPkgToInstall
	}
	pkgsToInstallStr := strings.Join(pkgsToInstall, ",")
	logger.Cmd.Info("installing " + pkgsToInstallStr)
	start := time.Now()
	cmd := exec.Command("apt", append([]string{"install", "-y"}, pkgsToInstall...)...)
	err := cmd.Run()
	since := time.Since(start)
	if err == nil {
		logger.Cmd.Info(pkgsToInstallStr+" installed", logger.Fields("duration", since.String()))
	} else {
		logger.Cmd.Error("failed to install "+pkgsToInstallStr, logger.Fields("duration", since.String()))
	}
	return cmd, err
}
