package adapters

import (
	"fmt"
	"os/exec"
	"strings"

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
	return runSystemCommand("npm", "install", pkgName)
}

func (n *NpmManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("npm", "uninstall", pkgName)
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

func runSystemCommand(name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func parseSearchLines(output, source string) []manager.SearchResult {
	var results []manager.SearchResult
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
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