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

func RestartService(services ...string) (*exec.Cmd, error) {
	if len(services) <= 0 {
		return nil, ErrNoServiceToRestart
	}
	servicesStr := strings.Join(services, ", ")
	logger.Cmd.Info("restarting " + servicesStr)
	start := time.Now()
	cmd := exec.Command("systemctl", append([]string{"restart"}, services...)...)
	err := cmd.Run()
	since := time.Since(start)
	if err == nil {
		logger.Cmd.Info(servicesStr+" restarted", logger.Fields("duration", since.String()))
	} else {
		logger.Cmd.Error("failed to restart "+servicesStr, logger.Fields("duration", since.String()))
	}
	return cmd, err
}
