package adapters

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// RpmManager manages packages on RPM-based systems.
type RpmManager struct{}

func (r *RpmManager) Name() string { return "rpm" }

func (r *RpmManager) IsAvailable() bool {
	_, err := exec.LookPath("rpm")
	return err == nil
}

func (r *RpmManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("rpm", "-i", pkgName)
}

func (r *RpmManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("rpm", "-e", pkgName)
}

func (r *RpmManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("rpm", "-qa").Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "rpm"), nil
}

func (r *RpmManager) Update() error {
	return runSystemCommand("rpm", "-Uvh", "--force")
}

func runSystemCommand(name string, args ...string) error {
	command := exec.Command("sudo", append([]string{name}, args...)...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func parseSearchLines(output, source string) []manager.SearchResult {
	var results []manager.SearchResult
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "Last metadata expiration") {
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