//go:build !windows
// +build !windows

package logger

import (
	"io"
	"log"
	"os"
)

// InitLogger sets up the logger to write to stdout, a log file, and syslog (on Unix).
func InitLogger(logFilePath string) error {
	// Create or open a log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Set up syslog
	/*
		sysLog, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "autopuller")
		if err != nil {
			return err
		}
	*/

	// Combine stdout, the log file, and syslog as the output destinations
	//multiWriter := io.MultiWriter(os.Stdout, logFile, sysLog)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// Set the log output to multiWriter
	log.SetOutput(multiWriter)

	// Set log flags (time, file, etc.)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Example log entry
	log.Println("Logger initialized (Unix)")

	return nil
}
