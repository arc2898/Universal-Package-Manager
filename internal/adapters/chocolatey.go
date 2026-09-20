package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// ChocolateyManager manages Windows packages using Chocolatey.
type ChocolateyManager struct{}

func (c *ChocolateyManager) Name() string { return "chocolatey" }

func (c *ChocolateyManager) IsAvailable() bool {
	_, err := exec.LookPath("choco")
	return err == nil
}

func (c *ChocolateyManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runWindowsCommand("choco", "install", "-y", pkgName)
}

func (c *ChocolateyManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runWindowsCommand("choco", "uninstall", "-y", pkgName)
}

func (c *ChocolateyManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("choco", "search", "--only-available", pkgName, "--limit-output", "10").Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "chocolatey"), nil
}

func (c *ChocolateyManager) Update() error {
	return runWindowsCommand("choco", "upgrade", "all", "-y")
}