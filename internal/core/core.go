package core

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/arc2898/Universal-Package-Manager/internal/adapters"
	"github.com/arc2898/Universal-Package-Manager/internal/detectors"
	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
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
	d.Register(&adapters.RpmManager{})
	d.Register(&adapters.HomebrewManager{})
	d.Register(&adapters.WingetManager{})
	d.Register(&adapters.ChocolateyManager{})
	d.Register(&adapters.ScoopManager{})
	d.Register(&adapters.CargoManager{})
	d.Register(&adapters.NpmManager{})
	d.Register(&adapters.PipManager{})
	d.Register(&adapters.GemManager{})
	d.Register(&adapters.XbpsManager{})
	d.Register(&adapters.ZypperManager{})
	d.Register(&adapters.PortageManager{})
	d.Register(&adapters.MasManager{})
	return &UPM{detector: d}
}

func (u *UPM) Init() { u.managers = u.detector.Detect() }
func (u *UPM) ListManagers() {
	fmt.Println("Available Package Managers:")
	for _, m := range u.managers { fmt.Printf("- %s\n", m.Name()) }
}
func (u *UPM) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil { return err }
	m, err := u.defaultManager()
	if err != nil { return err }
	return m.Install(pkgName)
}
func (u *UPM) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil { return err }
	m, err := u.defaultManager()
	if err != nil { return err }
	return m.Remove(pkgName)
}
func (u *UPM) Search(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil { return err }
	if len(u.managers) == 0 { return errors.New("no package managers detected") }
	var failures []error
	found := false
	fmt.Printf("Searching for '%s' across all managers...\n", pkgName)
	for _, m := range u.managers {
		results, err := m.Search(pkgName)
		if err != nil { failures = append(failures, fmt.Errorf("%s search: %w", m.Name(), err)); continue }
		if len(results) > 0 {
			found = true
			for _, result := range results { fmt.Printf("[%s] %s: %s\n", result.Source, result.Name, result.Description) }
		}
	}
	if !found { return errors.Join(failures...) }
	return nil
}
func (u *UPM) UpdateAll() error {
	if len(u.managers) == 0 { return errors.New("no package managers detected") }
	var failures []error
	for _, m := range u.managers {
		switch m.Name() {
		case "rpm", "cargo", "pip":
			continue
		}
		fmt.Printf("Updating %s...\n", m.Name())
		if err := m.Update(); err != nil { failures = append(failures, fmt.Errorf("%s update: %w", m.Name(), err)) }
	}
	return errors.Join(failures...)
}
func (u *UPM) BulkInstall(items map[string][]string) error { return u.bulk(items, "install") }
func (u *UPM) BulkRemove(items map[string][]string) error { return u.bulk(items, "remove") }
func (u *UPM) bulk(items map[string][]string, action string) error {
	var failures []error
	for managerName, packages := range items {
		m, err := u.managerByName(managerName)
		if err != nil { failures = append(failures, err); continue }
		for _, pkgName := range packages {
			if err := validatePackageName(pkgName); err != nil { failures = append(failures, fmt.Errorf("%s: %w", managerName, err)); continue }
			fmt.Printf("%sing %s via %s...\n", titleAction(action), pkgName, managerName)
			var operationErr error
			if action == "install" { operationErr = m.Install(pkgName) } else { operationErr = m.Remove(pkgName) }
			if operationErr != nil { failures = append(failures, fmt.Errorf("%s %s: %w", managerName, pkgName, operationErr)) }
		}
	}
	return errors.Join(failures...)
}
func (u *UPM) UninstallSelf() error {
	executable, err := os.Executable()
	if err != nil { return fmt.Errorf("locate executable: %w", err) }
	if err := os.Remove(executable); err == nil || os.IsNotExist(err) { return nil } else if runtime.GOOS == "windows" { return fmt.Errorf("remove %s after UPM exits: %w", executable, err) }
	if err := exec.Command("sudo", "rm", "-f", executable).Run(); err != nil { return fmt.Errorf("remove %s: %w", executable, err) }
	return nil
}
func (u *UPM) defaultManager() (manager.Manager, error) {
	if len(u.managers) == 0 { return nil, errors.New("no package managers detected") }
	for _, name := range nativeManagerPreference() {
		if m, err := u.managerByName(name); err == nil { return m, nil }
	}
	return u.managers[0], nil
}
func (u *UPM) managerByName(name string) (manager.Manager, error) {
	for _, m := range u.managers { if m.Name() == name { return m, nil } }
	return nil, fmt.Errorf("package manager %q is not detected", name)
}

func nativeManagerPreference() []string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil { return nil }
	return nativeManagerPreferenceFromOSRelease(string(data))
}

func nativeManagerPreferenceFromOSRelease(data string) []string {
	var id, like string
	for _, line := range strings.Split(data, "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok { continue }
		value = strings.Trim(value, "\"'")
		switch key {
		case "ID": id = strings.ToLower(value)
		case "ID_LIKE": like = strings.ToLower(value)
		}
	}
	switch id {
	case "debian", "ubuntu": return []string{"apt"}
	case "fedora", "rhel", "centos", "rocky", "almalinux": return []string{"dnf"}
	case "arch", "manjaro", "endeavouros": return []string{"pacman"}
	case "opensuse", "opensuse-leap", "opensuse-tumbleweed", "sles", "sled", "suse": return []string{"zypper"}
	case "void": return []string{"xbps"}
	case "gentoo": return []string{"portage"}
	}
	switch {
	case strings.Contains(like, "debian"): return []string{"apt"}
	case strings.Contains(like, "rhel") || strings.Contains(like, "fedora"): return []string{"dnf"}
	case strings.Contains(like, "arch"): return []string{"pacman"}
	case strings.Contains(like, "suse"): return []string{"zypper"}
	case strings.Contains(like, "void"): return []string{"xbps"}
	case strings.Contains(like, "gentoo"): return []string{"portage"}
	default: return nil
	}
}
func titleAction(action string) string {
	if action == "" { return "" }
	return strings.ToUpper(action[:1]) + action[1:]
}
func validatePackageName(pkgName string) error {
	if strings.TrimSpace(pkgName) == "" { return errors.New("package name cannot be empty") }
	if strings.HasPrefix(pkgName, "-") { return fmt.Errorf("package name cannot begin with '-': %q", pkgName) }
	if strings.IndexFunc(pkgName, func(r rune) bool { return r <= 0x20 || r == 0x7f }) >= 0 { return errors.New("package name cannot contain whitespace or control characters") }
	return nil
}
