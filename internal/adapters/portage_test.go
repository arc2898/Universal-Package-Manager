package adapters

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestPortageInstallDoesNotPrompt(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("fake sudo/emerge use POSIX shell scripts")
	}
	dir := t.TempDir()
	logPath := filepath.Join(dir, "emerge-args")
	emerge := filepath.Join(dir, "emerge")
	sudo := filepath.Join(dir, "sudo")
	emergeScript := "#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$LOG_FILE\"\n"
	sudoScript := "#!/bin/sh\nexec \"$@\"\n"
	if err := os.WriteFile(emerge, []byte(emergeScript), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(sudo, []byte(sudoScript), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("LOG_FILE", logPath)
	if err := (&PortageManager{}).Install("vim"); err != nil {
		t.Fatalf("Install returned error: %v", err)
	}
	args, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(args), "--ask=n\n") || !strings.HasSuffix(string(args), "vim\n") {
		t.Fatalf("expected non-interactive emerge args, got %q", args)
	}
}
