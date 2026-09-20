package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// GemManager manages Ruby packages using gem.
type GemManager struct{}

func (g *GemManager) Name() string { return "gem" }

func (g *GemManager) IsAvailable() bool {
	_, err := exec.LookPath("gem")
	return err == nil
}

func (g *GemManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("gem", "install", pkgName)
}

func (g *GemManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("gem", "uninstall", pkgName)
}

func (g *GemManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("gem", "search", pkgName, "--remote").Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLinesGem(string(output), "gem"), nil
}

func (g *GemManager) Update() error {
	return runSystemCommand("gem", "update")
}