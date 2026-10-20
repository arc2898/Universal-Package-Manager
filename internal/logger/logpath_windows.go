//go:build windows

package logger

import (
	"os"
	"path/filepath"
)

func logPath() string {
	if configured := os.Getenv("UPM_LOG_FILE"); configured != "" {
		return configured
	}
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		return filepath.Join(home, ".local", "log", "upm.log")
	}
	return "upm.log"
}
