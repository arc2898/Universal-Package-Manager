package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// ScoopManager manages Windows packages using Scoop.
type ScoopManager struct{}

func (s *ScoopManager) Name() string { return "scoop" }

func (s *ScoopManager) IsAvailable() bool {
	_, err := exec.LookPath("scoop")
	return err == nil
}

func (s *ScoopManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("scoop", "install", pkgName)
}

func (s *ScoopManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("scoop", "uninstall", pkgName)
}

func (s *ScoopManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("scoop", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "scoop"), nil
}

func (s *ScoopManager) Update() error {
	return runCommand("scoop", "update")
}