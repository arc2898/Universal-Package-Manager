package adapters

import (
	"os/exec"

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