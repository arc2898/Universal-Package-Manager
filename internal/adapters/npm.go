package adapters

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

// NpmManager manages Node.js packages using npm.
type NpmManager struct{}

func (n *NpmManager) Name() string { return "npm" }

func (n *NpmManager) IsAvailable() bool {
	_, err := exec.LookPath("npm")
	return err == nil
}

func (n *NpmManager) Install(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("npm", "install", "-g", pkgName)
}

func (n *NpmManager) Remove(pkgName string) error {
	if err := validatePackageName(pkgName); err != nil {
		return err
	}
	return runCommand("npm", "uninstall", "-g", pkgName)
}

func (n *NpmManager) Search(pkgName string) ([]manager.SearchResult, error) {
	if err := validatePackageName(pkgName); err != nil {
		return nil, err
	}
	output, err := exec.Command("npm", "search", pkgName, "--json").Output()
	if err != nil {
		return nil, err
	}
	var packages []struct {
		Name        string "json:\"name\""
		Description string "json:\"description\""
		Version     string "json:\"version\""
	}
	if err := json.Unmarshal(output, &packages); err != nil {
		return nil, fmt.Errorf("parse npm search results: %w", err)
	}
	results := make([]manager.SearchResult, 0, len(packages))
	for _, pkg := range packages {
		if pkg.Name == "" {
			continue
		}
		results = append(results, manager.SearchResult{
			Name: pkg.Name, Description: pkg.Description, Version: pkg.Version, Source: "npm",
		})
	}
	return results, nil
}

func (n *NpmManager) Update() error {
	return runCommand("npm", "update", "-g")
}