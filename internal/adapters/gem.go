package adapters

import (
	"fmt"
	"os/exec"
	"strings"

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
	return parseSearchLines(string(output), "gem"), nil
}

func (g *GemManager) Update() error {
	return runSystemCommand("gem", "update")
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
		if parts := strings.SplitN(line, " ", 2); len(parts) == 2 {
			name, description = parts[0], parts[1]
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