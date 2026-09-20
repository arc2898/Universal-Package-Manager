package adapters

import (
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// XbpsManager manages packages on Void Linux.
type XbpsManager struct{}

func (x *XbpsManager) Name() string { return "xbps" }

func (x *XbpsManager) IsAvailable() bool {
	_, err := exec.LookPath("xbps-install")
	return err == nil
}

func (x *XbpsManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("xbps-install", "-S", pkgName)
}

func (x *XbpsManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("xbps-remove", "-R", pkgName)
}

func (x *XbpsManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("xbps-query", "-Rs", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "xbps"), nil
}

func (x *XbpsManager) Update() error {
	return runSystemCommand("xbps-install", "-Su")
}

// ZypperManager manages packages on openSUSE/SUSE.
type ZypperManager struct{}

func (z *ZypperManager) Name() string { return "zypper" }

func (z *ZypperManager) IsAvailable() bool {
	_, err := exec.LookPath("zypper")
	return err == nil
}

func (z *ZypperManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("zypper", "install", "-y", pkgName)
}

func (z *ZypperManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("zypper", "remove", "-y", pkgName)
}

func (z *ZypperManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("zypper", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "zypper"), nil
}

func (z *ZypperManager) Update() error {
	return runSystemCommand("zypper", "refresh", "&&", "zypper", "update", "-y")
}

// PortageManager manages packages on Gentoo.
type PortageManager struct{}

func (p *PortageManager) Name() string { return "portage" }

func (p *PortageManager) IsAvailable() bool {
	_, err := exec.LookPath("emerge")
	return err == nil
}

func (p *PortageManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("emerge", "--ask", pkgName)
}

func (p *PortageManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runSystemCommand("emerge", "--unmerge", pkgName)
}

func (p *PortageManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("emerge", "-S", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "portage"), nil
}

func (p *PortageManager) Update() error {
	return runSystemCommand("emerge", "--sync", "&&", "emerge", "-uDN", "@world")
}

// MasManager manages Mac App Store packages using mas.
type MasManager struct{}

func (m *MasManager) Name() string { return "mas" }

func (m *MasManager) IsAvailable() bool {
	_, err := exec.LookPath("mas")
	return err == nil
}

func (m *MasManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("mas", "install", pkgName)
}

func (m *MasManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	// mas doesn't support uninstall directly, but we can try
	return runCommand("mas", "uninstall", pkgName)
}

func (m *MasManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("mas", "search", pkgName).Output()
	if err != nil {
		return nil, err
	}
	return parseSearchLines(string(output), "mas"), nil
}

func (m *MasManager) Update() error {
	return runCommand("mas", "upgrade")
}