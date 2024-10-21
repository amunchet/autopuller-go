//go:build windows
// +build windows

package logger

import (
	"io"
	"log"
	"os"
)

// InitLogger sets up the logger to write to stdout and a log file (without syslog on Windows).
func InitLogger() {
	// Create or open a log file
	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Combine stdout and the log file as the output destinations
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set the log output to multiWriter
	log.SetOutput(multiWriter)

	// Set log flags (time, file, etc.)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Example log entry
	log.Println("Logger initialized (Windows)")
}
