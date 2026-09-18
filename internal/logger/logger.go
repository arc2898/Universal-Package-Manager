package logger

import (
	"fmt"
	"os"
	"time"
)

const logFile = "/var/log/upm.log"

func Log(message string) error {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o640)
	if err != nil {
		return fmt.Errorf("open log: %w", err)
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	if _, err := fmt.Fprintf(file, "[%s] %s\n", timestamp, message); err != nil {
		return fmt.Errorf("write log: %w", err)
	}
	return nil
}

func ReadLogs() error {
	data, err := os.ReadFile(logFile)
	if os.IsNotExist(err) {
		fmt.Println("No logs found.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("read log: %w", err)
	}
	fmt.Print(string(data))
	return nil
}
