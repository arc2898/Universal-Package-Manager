package core

import (
	"bytes"
	"strings"
	"testing"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

func TestValidatePackageName(t *testing.T) {
	for _, name := range []string{"git", "@scope/package", "python3-pip", "foo.bar"} {
		if err := validatePackageName(name); err != nil {
			t.Errorf("validatePackageName(%q) unexpectedly failed: %v", name, err)
		}
	}
	for _, name := range []string{"", "  ", "--help", "foo bar", "foo\nbar", "foo\x00bar"} {
		if err := validatePackageName(name); err == nil {
			t.Errorf("validatePackageName(%q) expected an error", name)
		}
	}
}

type fakeManager struct {
	name      string
	results   []manager.SearchResult
	installed string
}

func (f *fakeManager) Name() string                                  { return f.name }
func (f *fakeManager) IsAvailable() bool                             { return true }
func (f *fakeManager) Install(pkg string) error                      { f.installed = pkg; return nil }
func (f *fakeManager) Remove(string) error                           { return nil }
func (f *fakeManager) Search(string) ([]manager.SearchResult, error) { return f.results, nil }
func (f *fakeManager) Update() error                                 { return nil }

func TestFindPackagesExactAndSimilar(t *testing.T) {
	apt := &fakeManager{name: "apt", results: []manager.SearchResult{
		{Name: "vlc", Source: "apt"},
		{Name: "vlc-plugin-base", Source: "apt"},
	}}
	flatpak := &fakeManager{name: "flatpak", results: []manager.SearchResult{
		{Name: "org.videolan.VLC", Aliases: []string{"VLC"}, Source: "flatpak"},
	}}
	u := &UPM{managers: []manager.Manager{apt, flatpak}}

	exact, err := u.FindPackages("vlc", false)
	if err != nil {
		t.Fatalf("FindPackages exact returned error: %v", err)
	}
	if len(exact) != 2 || exact[0].Manager.Name() != "apt" || exact[1].Manager.Name() != "flatpak" {
		t.Fatalf("unexpected exact matches: %#v", exact)
	}

	similar, err := u.FindPackages("vlc", true)
	if err != nil {
		t.Fatalf("FindPackages similar returned error: %v", err)
	}
	if len(similar) != 3 {
		t.Fatalf("got %d similar matches, want 3: %#v", len(similar), similar)
	}
}

func TestInstallFromPromptsAndInstallsSelectedSource(t *testing.T) {
	apt := &fakeManager{name: "apt", results: []manager.SearchResult{{Name: "vlc", Source: "apt"}}}
	snap := &fakeManager{name: "snap", results: []manager.SearchResult{{Name: "vlc", Source: "snap"}}}
	u := &UPM{managers: []manager.Manager{apt, snap}}

	var output bytes.Buffer
	if err := u.InstallFrom("vlc", false, strings.NewReader("2\n"), &output); err != nil {
		t.Fatalf("InstallFrom returned error: %v", err)
	}
	if apt.installed != "" || snap.installed != "vlc" {
		t.Fatalf("wrong manager installed package: apt=%q snap=%q", apt.installed, snap.installed)
	}
	if !strings.Contains(output.String(), "apt") || !strings.Contains(output.String(), "snap") {
		t.Fatalf("expected both sources in prompt, got %q", output.String())
	}
}

func TestInstallFromCanCancel(t *testing.T) {
	apt := &fakeManager{name: "apt", results: []manager.SearchResult{{Name: "vlc", Source: "apt"}}}
	u := &UPM{managers: []manager.Manager{apt}}
	if err := u.InstallFrom("vlc", false, strings.NewReader("0\n"), &bytes.Buffer{}); err != nil {
		t.Fatalf("cancelling install should not fail: %v", err)
	}
	if apt.installed != "" {
		t.Fatalf("package was installed after cancellation: %q", apt.installed)
	}
}
