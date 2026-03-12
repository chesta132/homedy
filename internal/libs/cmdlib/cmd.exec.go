package cmdlib

import (
	"homedy/internal/libs/logger"
	"os/exec"
	"strings"
)

func InstallPkgs(pkgs []string) *exec.Cmd {
	pkgsToInstall := filterUninstalledPkgs(pkgs)
	if len(pkgsToInstall) <= 0 {
		return nil
	}
	logger.Cmd.Info("installing " + strings.Join(pkgsToInstall, ","))
	pkgsToInstall = append([]string{"install"}, pkgsToInstall...)
	return exec.Command("apt", pkgsToInstall...)
}
