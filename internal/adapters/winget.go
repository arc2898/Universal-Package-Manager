package adapters

import (
	"os/exec"

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
	return runWindowsCommand("winget", "install", "--id", pkgName, "--silent")
}

func (w *WingetManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runWindowsCommand("winget", "uninstall", "--id", pkgName)
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
	return runWindowsCommand("winget", "upgrade", "--all", "--silent")
}