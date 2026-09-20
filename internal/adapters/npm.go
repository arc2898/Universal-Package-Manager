package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// NpmManager manages Node.js packages using npm.
type NpmManager struct{}

func (n *NpmManager) Name() string { return "npm" }

func (n *NpmManager) IsAvailable() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

func (n *NpmManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("npm", "install", "-g", pkgName)
}

func (n *NpmManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("npm", "uninstall", "-g", pkgName)
}

func (n *NpmManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("npm", "search", pkgName, "--json").Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "npm"), nil
}

func (n *NpmManager) Update() error {
	return runSystemCommand("npm", "update", "-g")
}