package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// PipManager manages Python packages using pip.
type PipManager struct{}

func (p *PipManager) Name() string { return "pip" }

func (p *PipManager) IsAvailable() bool {
	_, err := exec.LookPath("pip")
	return err == nil
}

func (p *PipManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("pip", "install", pkgName)
}

func (p *PipManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("pip", "uninstall", "-y", pkgName)
}

func (p *PipManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	// "pip search" was removed from modern pip/PyPI. "pip index versions"
	// supports an exact package lookup against the configured package index.
	output, err := exec.Command("pip", "index", "versions", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parsePipIndexVersions(string(output), "pip"), nil
}

func (p *PipManager) Update() error {
	return runCommand("pip", "install", "--upgrade", "pip")
}