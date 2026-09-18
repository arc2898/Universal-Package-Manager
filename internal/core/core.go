package core

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/manus/upm/internal/adapters"
	"github.com/manus/upm/internal/detectors"
	"github.com/manus/upm/pkg/manager"
)

type UPM struct {
	detector *detectors.Detector
	managers []manager.Manager
}

func NewUPM() *UPM {
	d := detectors.NewDetector()
	d.Register(&adapters.AptManager{})
	d.Register(&adapters.SnapManager{})
	d.Register(&adapters.FlatpakManager{})
	d.Register(&adapters.DnfManager{})
	d.Register(&adapters.PacmanManager{})
	return &UPM{detector: d}
}

func (u *UPM) Init() { u.managers = u.detector.Detect() }

func (u *UPM) ListManagers() {
	fmt.Println("Available Package Managers:")
	for _, m := range u.managers {
		fmt.Printf("- %s\n", m.Name())
	}
}

func (u *UPM) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	m, err := u.defaultManager()
	if err != nil {
		return err
	}
	return m.Install(pkgName)
}

func (u *UPM) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	m, err := u.defaultManager()
	if err != nil {
		return err
	}
	return m.Remove(pkgName)
}

func (u *UPM) Search(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	if len(u.managers) == 0 {
		return errors.New("no package managers detected")
	}

	var failures []error
	fmt.Printf("Searching for '%s' across all managers...\n", pkgName)
	for _, m := range u.managers {
		results, err := m.Search(pkgName)
		if err != nil {
			failures = append(failures, fmt.Errorf("%s search: %w", m.Name(), err))
			continue
		}
		for _, result := range results {
			fmt.Printf("[%s] %s: %s\n", result.Source, result.Name, result.Description)
		}
	}
	return errors.Join(failures...)
}

func (u *UPM) UpdateAll() error {
	if len(u.managers) == 0 {
		return errors.New("no package managers detected")
	}
	var failures []error
	for _, m := range u.managers {
		fmt.Printf("Updating %s...\n", m.Name())
		if err := m.Update(); err != nil {
			failures = append(failures, fmt.Errorf("%s update: %w", m.Name(), err))
		}
	}
	return errors.Join(failures...)
}

func (u *UPM) BulkInstall(items map[string][]string) error {
	return u.bulk(items, "install")
}

func (u *UPM) BulkRemove(items map[string][]string) error {
	return u.bulk(items, "remove")
}

func (u *UPM) bulk(items map[string][]string, action string) error {
	var failures []error
	for managerName, packages := range items {
		m, err := u.managerByName(managerName)
		if err != nil {
			failures = append(failures, err)
			continue
		}
		for _, pkgName := range packages {
			if err := validatePackageName(pkgName); err != nil {
				failures = append(failures, fmt.Errorf("%s: %w", managerName, err))
				continue
			}
			fmt.Printf("%sing %s via %s...\n", strings.Title(action), pkgName, managerName)
			var operationErr error
			if action == "install" {
				operationErr = m.Install(pkgName)
			} else {
				operationErr = m.Remove(pkgName)
			}
			if operationErr != nil {
				failures = append(failures, fmt.Errorf("%s %s: %w", managerName, pkgName, operationErr))
			}
		}
	}
	return errors.Join(failures...)
}

func (u *UPM) UninstallSelf() error {
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("locate executable: %w", err)
	}
	if err := exec.Command("sudo", "rm", "-f", executable).Run(); err != nil {
		return fmt.Errorf("remove %s: %w", executable, err)
	}
	return nil
}

func (u *UPM) defaultManager() (manager.Manager, error) {
	if len(u.managers) == 0 {
		return nil, errors.New("no package managers detected")
	}
	preferred := nativeManagerPreference()
	for _, name := range preferred {
		if m, err := u.managerByName(name); err == nil {
			return m, nil
		}
	}
	return u.managers[0], nil
}

func (u *UPM) managerByName(name string) (manager.Manager, error) {
	for _, m := range u.managers {
		if m.Name() == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("package manager %q is not detected", name)
}

func nativeManagerPreference() []string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return nil
	}
	values := string(data)
	switch {
	case strings.Contains(values, "ID=ubuntu") || strings.Contains(values, "ID=debian") || strings.Contains(values, "ID_LIKE=debian"):
		return []string{"apt"}
	case strings.Contains(values, "ID=fedora") || strings.Contains(values, "ID=rhel") || strings.Contains(values, "ID=centos") || strings.Contains(values, "ID_LIKE="+"rhel"):
		return []string{"dnf"}
	case strings.Contains(values, "ID=arch") || strings.Contains(values, "ID=manjaro") || strings.Contains(values, "ID_LIKE=arch"):
		return []string{"pacman"}
	default:
		return nil
	}
}

func validatePackageName(pkgName string) error {
	if strings.TrimSpace(pkgName) == "" {
		return errors.New("package name cannot be empty")
	}
	return nil
}
