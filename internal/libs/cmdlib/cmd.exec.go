package cmdlib

import (
	"homedy/internal/libs/logger"
	"os/exec"
	"strings"
)

func InstallPkgs(pkgs []string) (*exec.Cmd, error) {
	pkgsToInstall := filterUninstalledPkgs(pkgs)
	if len(pkgsToInstall) <= 0 {
		return nil, ErrNoPkgToInstall
	}
	pkgsToInstallStr := strings.Join(pkgsToInstall, ",")
	logger.Cmd.Info("installing " + pkgsToInstallStr)
	pkgsToInstall = append([]string{"install"}, pkgsToInstall...)
	cmd := exec.Command("apt", pkgsToInstall...)
	err := cmd.Run()
	if err != nil {
		logger.Cmd.Info(pkgsToInstallStr+" installed", logger.Fields("duration", cmd.WaitDelay.String()))
	}
	return cmd, err
}
