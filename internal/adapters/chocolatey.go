package adapters

import (
	"fmt"
	"os/exec"
	"strings"

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
	return runSystemCommand("choco", "install", "-y", pkgName)
}

func (c *ChocolateyManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("choco", "uninstall", "-y", pkgName)
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
	return runSystemCommand("choco", "upgrade", "all", "-y")
}

func runSystemCommand(name string, args ...string) error {
	command := exec.Command("cmd", append([]string{"/c", name}, args...)...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func parseSearchLines(output, source string) []manager.SearchResult {
	var results []manager.SearchResult
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "ERROR") || strings.HasPrefix(line, "Warning") {
			continue
		}
		name, description := line, ""
		if parts := strings.SplitN(line, " - ", 2); len(parts) == 2 {
			name, description = strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		} else if fields := strings.Fields(line); len(fields) > 0 {
			name = fields[0]
		}
		results = append(results, manager.SearchResult{Name: name, Description: description, Source: source})
	}
	return results
}

func validatePackageName(pkgName string) error {
	if strings.TrimSpace(pkgName) == "" {
		return fmt.Errorf("package name cannot be empty")
	}
	return nil
}