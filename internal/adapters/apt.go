package adapters

import (
	"os/exec"
	"strings"

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

	var results []manager.SearchResult
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " - ", 2)
		result := manager.SearchResult{Name: strings.TrimSpace(parts[0]), Source: "apt"}
		if len(parts) == 2 {
			result.Description = strings.TrimSpace(parts[1])
		}
		results = append(results, result)
	}
	return results, nil
}

func (a *AptManager) Update() error {
	return runSystemCommand("apt", "update")
}
