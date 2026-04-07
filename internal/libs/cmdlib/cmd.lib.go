package cmdlib

import (
	"homedy/internal/libs/slicelib"
	"os/exec"
	"strings"
)

func extractPkgName(pkg string) string {
	// "samba=2:4.19.5" -> "samba"
	if i := strings.Index(pkg, "="); i != -1 {
		return pkg[:i]
	}
	return pkg
}

func isPkgInstalled(pkg string) bool {
	cmd := exec.Command("dpkg", "-s", extractPkgName(pkg))
	return cmd.Run() == nil
}

func filterUninstalledPkgs(pkgs []string) []string {
	return slicelib.Filter(pkgs, func(idx int, pkg string) bool { return !isPkgInstalled(pkg) })
}

func filterInstalledPkgs(pkgs []string) []string {
	return slicelib.Filter(pkgs, func(idx int, pkg string) bool { return isPkgInstalled(pkg) })
}

func holdPkgVer(pkg string) error {
	return exec.Command("apt-mark", "hold", extractPkgName(pkg)).Run()
}
