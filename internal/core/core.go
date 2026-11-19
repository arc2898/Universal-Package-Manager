package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"runtime"
	"sort"
	"strings"
	"sync"

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
	// New adapters for parity with Python version
	d.Register(&adapters.XbpsManager{})
	d.Register(&adapters.ZypperManager{})
	d.Register(&adapters.PortageManager{})
	d.Register(&adapters.MasManager{})
	return &UPM{detector: d}
}

func (u *UPM) Init() { u.managers = u.detector.Detect() }

func (u *UPM) ListManagers() {
	fmt.Println("Available Package Managers:")
	for _, m := range u.managers {
		fmt.Printf("- %s\n", m.Name())
	}
}

type PackageMatch struct {
	Manager manager.Manager
	Result  manager.SearchResult
}

// FindPackages searches every detected manager concurrently. Exact mode only
// returns package names that match the user's query; similar mode returns the
// managers' broader search results.
func (u *UPM) FindPackages(pkgName string, similar bool) ([]PackageMatch, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	if len(u.managers) == 0 {
		return nil, errors.New("no package managers detected")
	}

	type response struct {
		manager manager.Manager
		results []manager.SearchResult
		err     error
	}
	responses := make(chan response, len(u.managers))
	var wg sync.WaitGroup
	for _, m := range u.managers {
		wg.Add(1)
		go func(m manager.Manager) {
			defer wg.Done()
			results, err := m.Search(pkgName)
			responses <- response{manager: m, results: results, err: err}
		}(m)
	}
	wg.Wait()
	close(responses)

	var matches []PackageMatch
	var failures []error
	for response := range responses {
		if response.err != nil {
			failures = append(failures, fmt.Errorf("%s search: %w", response.manager.Name(), response.err))
			continue
		}
		for _, result := range response.results {
			if result.Name == "" {
				continue
			}
			exactMatch := strings.EqualFold(result.Name, pkgName) || strings.EqualFold(path.Base(result.Name), pkgName)
			for _, alias := range result.Aliases {
				if strings.EqualFold(alias, pkgName) {
					exactMatch = true
					break
				}
			}
			if similar || exactMatch {
				if result.Source == "" {
					result.Source = response.manager.Name()
				}
				matches = append(matches, PackageMatch{Manager: response.manager, Result: result})
			}
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Manager.Name() == matches[j].Manager.Name() {
			return strings.ToLower(matches[i].Result.Name) < strings.ToLower(matches[j].Result.Name)
		}
		return matches[i].Manager.Name() < matches[j].Manager.Name()
	})
	return matches, errors.Join(failures...)
}

// Install searches all detected managers and installs only after the user
// explicitly selects a result. Similar mode is enabled by `upm install -s`.
func (u *UPM) Install(pkgName string, similar bool) error {
	return u.InstallFrom(pkgName, similar, os.Stdin, os.Stdout)
}

func (u *UPM) InstallFrom(pkgName string, similar bool, input io.Reader, output io.Writer) error {
	matches, err := u.FindPackages(pkgName, similar)
	if err != nil && len(matches) == 0 {
		if similar {
			return fmt.Errorf("no similar packages found for %q: %w", pkgName, err)
		}
		return fmt.Errorf("no exact package named %q found; try `upm install -s %s` for similar names (search errors: %w)", pkgName, pkgName, err)
	}
	if len(matches) == 0 {
		return fmt.Errorf("no package named %q found; try `upm install -s %s`", pkgName, pkgName)
	}
	if err != nil {
		fmt.Fprintf(output, "Warning: some package managers could not be searched: %v\n", err)
	}

	if similar {
		fmt.Fprintf(output, "Similar package matches for %q:\n", pkgName)
	} else {
		fmt.Fprintf(output, "The requested package %q is available from:\n", pkgName)
	}
	for i, match := range matches {
		description := strings.TrimSpace(match.Result.Description)
		if description != "" {
			fmt.Fprintf(output, "  %d) %s — %s (%s)\n", i+1, match.Result.Name, description, match.Manager.Name())
		} else {
			fmt.Fprintf(output, "  %d) %s (%s)\n", i+1, match.Result.Name, match.Manager.Name())
		}
	}
	fmt.Fprint(output, "Select a source number to install, or 0 to cancel: ")

	reader := bufio.NewReader(input)
	var selection int
	if _, err := fmt.Fscan(reader, &selection); err != nil {
		return fmt.Errorf("read package selection: %w", err)
	}
	if selection == 0 {
		fmt.Fprintln(output, "Installation cancelled.")
		return nil
	}
	if selection < 1 || selection > len(matches) {
		return fmt.Errorf("invalid selection %d; choose 1-%d or 0 to cancel", selection, len(matches))
	}

	selected := matches[selection-1]
	fmt.Fprintf(output, "Installing %s from %s...\n", selected.Result.Name, selected.Manager.Name())
	return selected.Manager.Install(selected.Result.Name)
}

// SetupManager performs first-time setup for managers that require repository
// initialization before package discovery can work.
func (u *UPM) SetupManager(name string) error {
	m, err := u.managerByName(name)
	if err != nil {
		return err
	}
	if name == "portage" {
		portage, ok := m.(*adapters.PortageManager)
		if !ok {
			return errors.New("Portage adapter is unavailable")
		}
		return portage.Setup()
	}
	return fmt.Errorf("setup is not required or supported for %q", name)
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
	matches, err := u.FindPackages(pkgName, true)
	fmt.Printf("Search results for %q:\n", pkgName)
	for _, match := range matches {
		fmt.Printf("[%s] %s: %s\n", match.Manager.Name(), match.Result.Name, match.Result.Description)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: some package managers could not be searched: %v\n", err)
	}
	if len(matches) == 0 && err != nil {
		return err
	}
	if len(matches) == 0 {
		return fmt.Errorf("no packages found for %q", pkgName)
	}
	return nil
}

func (u *UPM) UpdateAll() error {
	if len(u.managers) == 0 {
		return errors.New("no package managers detected")
	}
	var failures []error
	for _, m := range u.managers {
		switch m.Name() {
		case "rpm":
			fmt.Println("Skipping rpm: direct RPM has no repository-wide update operation; use the distribution's package manager.")
			continue
		case "cargo":
			fmt.Println("Skipping cargo: 'cargo update' changes dependencies in the current project, not globally installed packages.")
			continue
		case "pip":
			fmt.Println("Skipping pip: pip has no safe, standard command to upgrade every installed package.")
			continue
		}
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
			fmt.Printf("%sing %s via %s...\n", titleAction(action), pkgName, managerName)
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
	if err := os.Remove(executable); err == nil || os.IsNotExist(err) {
		return nil
	} else if runtime.GOOS == "windows" {
		return fmt.Errorf("remove %s after UPM exits: %w", executable, err)
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

func titleAction(action string) string {
	if action == "" {
		return ""
	}
	return strings.ToUpper(action[:1]) + action[1:]
}

func validatePackageName(pkgName string) error {
	if strings.TrimSpace(pkgName) == "" {
		return errors.New("package name cannot be empty")
	}
	if strings.HasPrefix(pkgName, "-") {
		return fmt.Errorf("package name cannot begin with '-': %q", pkgName)
	}
	if strings.IndexFunc(pkgName, func(r rune) bool { return r <= 0x20 || r == 0x7f }) >= 0 {
		return errors.New("package name cannot contain whitespace or control characters")
	}
	return nil
}
