package adapters

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// DnfManager manages packages on Fedora, RHEL, CentOS, and other DNF systems.
type DnfManager struct{}

func (d *DnfManager) Name() string { return "dnf" }

func (d *DnfManager) IsAvailable() bool {
	_, err := exec.LookPath("dnf")
	return err == nil
}

func (d *DnfManager) Install(pkgName string) error {
	return runSystemCommand("dnf", "install", "-y", pkgName)
}

func (d *DnfManager) Remove(pkgName string) error {
	return runSystemCommand("dnf", "remove", "-y", pkgName)
}

func (d *DnfManager) Search(pkgName string) ([]manager.SearchResult, error) {
	output, err := exec.Command("dnf", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "dnf"), nil
}

func (d *DnfManager) Update() error {
	return runSystemCommand("dnf", "upgrade", "-y")
}

// PacmanManager manages packages on Arch Linux and derivatives.
type PacmanManager struct{}

func (p *PacmanManager) Name() string { return "pacman" }

func (p *PacmanManager) IsAvailable() bool {
	_, err := exec.LookPath("pacman")
	return err == nil
}

func (p *PacmanManager) Install(pkgName string) error {
	return runSystemCommand("pacman", "-S", "--noconfirm", pkgName)
}

func (p *PacmanManager) Remove(pkgName string) error {
	return runSystemCommand("pacman", "-R", "--noconfirm", pkgName)
}

func (p *PacmanManager) Search(pkgName string) ([]manager.SearchResult, error) {
	output, err := exec.Command("pacman", "-Ss", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "pacman"), nil
}

func (p *PacmanManager) Update() error {
	// Arch does not support partial upgrades; synchronize and upgrade together.
	return runSystemCommand("pacman", "-Syu", "--noconfirm")
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
