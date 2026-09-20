package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// AptManager manages Debian and Ubuntu packages.
type AptManager struct{}

func (a *AptManager) Name() string { return "apt" }

func (a *AptManager) IsAvailable() bool {
	_, err := exec.LookPath("apt")
	return err == nil
}

func (a *AptManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("apt", "install", "-y", pkgName)
}

func (a *AptManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("apt", "remove", "-y", pkgName)
}

func (a *AptManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("apt-cache", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "apt"), nil
}

func (a *AptManager) Update() error {
	return runSystemCommand("apt", "update")
}