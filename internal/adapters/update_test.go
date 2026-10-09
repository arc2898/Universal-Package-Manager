package adapters

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestZypperManagerUpdateRunsSeparateCommands(t *testing.T) {
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
