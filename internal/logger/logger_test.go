package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogUsesConfiguredPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "upm.log")
	t.Setenv("UPM_LOG_FILE", path)
	if err := Log("test message"); err != nil {
		t.Fatalf("Log returned error: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(data), "test message") {
		t.Fatalf("log does not contain message: %q", data)
	}
}
