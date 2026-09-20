package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// HomebrewManager manages Homebrew packages on macOS.
type HomebrewManager struct{}

func (h *HomebrewManager) Name() string { return "homebrew" }

func (h *HomebrewManager) IsAvailable() bool {
	_, err := exec.LookPath("brew")
	return err == nil
}

func (h *HomebrewManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("brew", "install", pkgName)
}

func (h *HomebrewManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("brew", "remove", pkgName)
}

func (h *HomebrewManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("brew", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "homebrew"), nil
}

func (h *HomebrewManager) Update() error {
	return runCommand("brew", "update")
}