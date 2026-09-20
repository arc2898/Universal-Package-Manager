package adapters

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// runSystemCommand runs a command with sudo
func runSystemCommand(name string, args ...string) error {
	command := exec.Command("sudo", append([]string{name}, args...)...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

// runCommand runs a command without sudo
func runCommand(name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

// runWindowsCommand runs a command via cmd.exe on Windows
func runWindowsCommand(name string, args ...string) error {
	command := exec.Command("cmd", append([]string{"/c", name}, args...)...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

// parseSearchLines parses generic "name - description" or space-separated output
func parseSearchLines(output, source string) []manager.SearchResult {
	var results []manager.SearchResult
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Skip common header/error lines
		if strings.HasPrefix(line, "Last metadata expiration") ||
			strings.HasPrefix(line, "Error") ||
			strings.HasPrefix(line, "Warning") ||
			strings.HasPrefix(line, "ERROR") {
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

// parseSearchLinesPip parses pip's "name (version) - description" format
func parseSearchLinesPip(output, source string) []manager.SearchResult {
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

// parseSearchLinesGem parses gem's "name version" format
func parseSearchLinesGem(output, source string) []manager.SearchResult {
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