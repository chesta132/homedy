package cmdlib

import (
	"homedy/internal/libs/slicelib"
	"os/exec"
)

func isPkgInstalled(pkg string) bool {
	cmd := exec.Command("dpkg", "-s", pkg)
	return cmd.ProcessState.Success()
}

func filterUninstalledPkgs(pkgs []string) []string {
	return slicelib.Filter(pkgs, func(idx int, pkg string) bool { return isPkgInstalled(pkg) })
}

func filterInstalledPkgs(pkgs []string) []string {
	return slicelib.Filter(pkgs, func(idx int, pkg string) bool { return !isPkgInstalled(pkg) })
}
