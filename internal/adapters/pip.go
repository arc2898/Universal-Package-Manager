package adapters

import (
	"fmt"
	"os/exec"
	"strings"

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
	return parseSearchLines(string(output), "pip"), nil
}

func (p *PipManager) Update() error {
	return runSystemCommand("pip", "install", "--upgrade", "pip")
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
		if parts := strings.SplitN(line, " (", 2); len(parts) == 2 {
			name = strings.Trim(parts[0], " ")
			description = strings.Trim(parts[1], ")")
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