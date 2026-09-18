package adapters

import (
	"os"
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// SnapManager manages Snap packages.
type SnapManager struct{}

func (s *SnapManager) Name() string { return "snap" }

func (s *SnapManager) IsAvailable() bool {
	_, err := exec.LookPath("snap")
	return err == nil
}

func (s *SnapManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("snap", "install", pkgName)
}

func (s *SnapManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("snap", "remove", pkgName)
}

func (s *SnapManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("snap", "find", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "snap"), nil
}

func (s *SnapManager) Update() error {
	return runSystemCommand("snap", "refresh")
}

// FlatpakManager manages Flatpak applications and runtimes.
type FlatpakManager struct{}

func (f *FlatpakManager) Name() string { return "flatpak" }

func (f *FlatpakManager) IsAvailable() bool {
	_, err := exec.LookPath("flatpak")
	return err == nil
}

func (f *FlatpakManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("flatpak", "install", "-y", pkgName)
}

func (f *FlatpakManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("flatpak", "uninstall", "-y", pkgName)
}

func (f *FlatpakManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("flatpak", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "flatpak"), nil
}

func (f *FlatpakManager) Update() error {
	return runCommand("flatpak", "update", "-y")
}

func runCommand(name string, args ...string) error {
	command := exec.Command(name, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}
