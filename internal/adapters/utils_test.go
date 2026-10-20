package adapters

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/arc2898/Universal-Package-Manager/pkg/manager"
)

func TestValidatePackageName(t *testing.T) {
	valid := []string{"git", "@scope/package", "python3-pip", "foo.bar", "foo+bar"}
	for _, name := range valid {
		if err := validatePackageName(name); err != nil {
			t.Errorf("validatePackageName(%q) unexpectedly failed: %v", name, err)
		}
	}
	invalid := []string{"", "   ", "-y", "--help", "foo bar", "foo\nbar", "foo\x00bar"}
	for _, name := range invalid {
		if err := validatePackageName(name); err == nil {
			t.Errorf("validatePackageName(%q) expected an error", name)
		}
	}
}

func TestParsePipIndexVersions(t *testing.T) {
	got := parsePipIndexVersions("requests (2.32.3)\nAvailable versions: 2.32.3, 2.32.2\n", "pip")
	if len(got) != 1 {
		t.Fatalf("got %d results, want 1: %#v", len(got), got)
	}
	if got[0].Name != "requests" || got[0].Version != "2.32.3" || got[0].Source != "pip" {
		t.Fatalf("unexpected result: %#v", got[0])
	}
}

func TestNpmSearchParsesJSON(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fake npm executable in this test uses a POSIX shell script")
	}
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "npm-args")
	npmPath := filepath.Join(tempDir, "npm")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$LOG_FILE\"\nprintf '%s' '[{\"name\":\"typescript\",\"description\":\"TypeScript language tools\",\"version\":\"5.9.0\"}]'\n"
	if err := os.WriteFile(npmPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", logPath)
	results, err := (&NpmManager{}).Search("typescript")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 || results[0].Name != "typescript" || results[0].Version != "5.9.0" {
		t.Fatalf("unexpected npm search results: %#v", results)
	}
	args, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(args), "--json") {
		t.Fatalf("expected npm search to request JSON, args were %q", args)
	}
}

func TestRpmUpdateIsExplicitlyUnsupported(t *testing.T) {
	err := (&RpmManager{}).Update()
	if err == nil || !strings.Contains(err.Error(), "repository-wide updates") {
		t.Fatalf("expected safe unsupported-operation error, got %v", err)
	}
}

func TestPipSearchUsesSupportedIndexCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fake pip executable in this test uses a POSIX shell script")
	}
	tempDir := t.TempDir()
	argsPath := filepath.Join(tempDir, "pip-args")
	pipPath := filepath.Join(tempDir, "pip")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$LOG_FILE\"\nprintf '%s\\n' 'requests (2.32.3)' 'Available versions: 2.32.3, 2.32.2'\n"
	if err := os.WriteFile(pipPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", argsPath)
	results, err := (&PipManager{}).Search("requests")
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 || results[0].Name != "requests" || results[0].Version != "2.32.3" {
		t.Fatalf("unexpected pip search results: %#v", results)
	}
	args, err := os.ReadFile(argsPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(args), "index\nversions\nrequests") {
		t.Fatalf("expected pip index versions command, args were %q", args)
	}
}

func TestUnsupportedOptionLikePackagesRejectedBeforeCommands(t *testing.T) {
	managers := []struct {
		name string
		manager interface {
			Install(string) error
			Remove(string) error
			Search(string) ([]manager.SearchResult, error)
		}
	}{
		{name: "dnf", manager: &DnfManager{}},
		{name: "pacman", manager: &PacmanManager{}},
		{name: "rpm", manager: &RpmManager{}},
	}
	for _, item := range managers {
		t.Run(item.name, func(t *testing.T) {
			if err := item.manager.Install("--help"); err == nil {
				t.Fatal("Install accepted an option-like package name")
			}
			if err := item.manager.Remove("--help"); err == nil {
				t.Fatal("Remove accepted an option-like package name")
			}
			if _, err := item.manager.Search("--help"); err == nil {
				t.Fatal("Search accepted an option-like package name")
			}
		})
	}
}

func TestAptUpdateRefreshesAndUpgrades(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fake sudo executable in this test uses a POSIX shell script")
	}
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "commands.log")
	sudoPath := filepath.Join(tempDir, "sudo")
	if err := os.WriteFile(sudoPath, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$LOG_FILE\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", logPath)
	if err := (&AptManager{}).Update(); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	log := string(data)
	refresh := strings.Index(log, "apt\nupdate\n")
	upgrade := strings.Index(log, "apt\nupgrade\n-y\n")
	if refresh < 0 || upgrade < 0 || refresh >= upgrade {
		t.Fatalf("expected apt update before apt upgrade, got %q", log)
	}
}
