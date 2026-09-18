package adapters

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// WingetManager manages Windows packages using winget.
type WingetManager struct{}

func (w *WingetManager) Name() string { return "winget" }

func (w *WingetManager) IsAvailable() bool {
	_, err := exec.LookPath("winget")
	return err == nil
}

func (w *WingetManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("winget", "install", "--id", pkgName, "--silent")
}

func (w *WingetManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("winget", "uninstall", "--id", pkgName)
}

func (w *WingetManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("winget", "search", "--eula", "off", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "winget"), nil
}

func (w *WingetManager) Update() error {
	return runSystemCommand("winget", "upgrade", "--all", "--silent")
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