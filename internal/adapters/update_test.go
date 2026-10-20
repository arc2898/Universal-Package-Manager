package adapters

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestZypperManagerUpdateRunsSeparateCommands(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fake sudo executable in this test uses a POSIX shell script")
	}
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "sudo.log")
	sudoPath := filepath.Join(tempDir, "sudo")
	if err := os.WriteFile(sudoPath, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$LOG_FILE\"\n"), 0o755); err != nil {
		t.Fatalf("write sudo stub: %v", err)
	}
	t.Setenv("PATH", tempDir+":"+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", logPath)

	if err := (&ZypperManager{}).Update(); err != nil {
		t.Fatalf("ZypperManager.Update() returned error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read sudo log: %v", err)
	}
	log := string(data)
	if strings.Contains(log, "&&") {
		t.Fatalf("expected chained shell operators to be split into separate invocations, got: %q", log)
	}
	if !strings.Contains(log, "zypper\nrefresh\n") || !strings.Contains(log, "zypper\nupdate\n-y\n") {
		t.Fatalf("expected refresh and update to run as separate commands, got: %q", log)
	}
}

func TestPortageManagerUpdateRunsSeparateCommands(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fake sudo executable in this test uses a POSIX shell script")
	}
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "sudo.log")
	sudoPath := filepath.Join(tempDir, "sudo")
	if err := os.WriteFile(sudoPath, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$LOG_FILE\"\n"), 0o755); err != nil {
		t.Fatalf("write sudo stub: %v", err)
	}
	t.Setenv("PATH", tempDir+":"+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", logPath)

	if err := (&PortageManager{}).Update(); err != nil {
		t.Fatalf("PortageManager.Update() returned error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read sudo log: %v", err)
	}
	log := string(data)
	if strings.Contains(log, "&&") {
		t.Fatalf("expected chained shell operators to be split into separate invocations, got: %q", log)
	}
	if !strings.Contains(log, "emerge\n--sync\n") || !strings.Contains(log, "emerge\n-uDN\n@world\n") {
		t.Fatalf("expected sync and world upgrade to run as separate commands, got: %q", log)
	}
}


func TestHomebrewUpdateRefreshesAndUpgrades(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fake brew executable in this test uses a POSIX shell script")
	}
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "brew.log")
	brewPath := filepath.Join(tempDir, "brew")
	if err := os.WriteFile(brewPath, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$LOG_FILE\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", logPath)
	if err := (&HomebrewManager{}).Update(); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	log := string(data)
	refresh := strings.Index(log, "update\n")
	upgrade := strings.Index(log, "upgrade\n")
	if refresh < 0 || upgrade < 0 || refresh >= upgrade {
		t.Fatalf("expected brew update before brew upgrade, got %q", log)
	}
}

func TestScoopUpdateIncludesAllApps(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the fake scoop executable in this test uses a POSIX shell script")
	}
	tempDir := t.TempDir()
	logPath := filepath.Join(tempDir, "scoop.log")
	scoopPath := filepath.Join(tempDir, "scoop")
	if err := os.WriteFile(scoopPath, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$LOG_FILE\"\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", logPath)
	if err := (&ScoopManager{}).Update(); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "update\n*\n") {
		t.Fatalf("expected Scoop to update all apps, got %q", data)
	}
}
