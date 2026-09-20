package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// CargoManager manages Rust packages using cargo.
type CargoManager struct{}

func (c *CargoManager) Name() string { return "cargo" }

func (c *CargoManager) IsAvailable() bool {
	_, err := exec.LookPath("cargo")
	return err == nil
}

func (c *CargoManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("cargo", "install", pkgName)
}

func (c *CargoManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("cargo", "uninstall", pkgName)
}

func (c *CargoManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("cargo", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "cargo"), nil
}

func (c *CargoManager) Update() error {
	return runSystemCommand("cargo", "update")
}