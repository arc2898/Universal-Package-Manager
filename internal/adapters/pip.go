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
	return runSystemCommand("pip", "install", pkgName)
}

func (p *PipManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("pip", "uninstall", "-y", pkgName)
}

func (p *PipManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("pip", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLinesPip(string(output), "pip"), nil
}

func (p *PipManager) Update() error {
	return runSystemCommand("pip", "install", "--upgrade", "pip")
}