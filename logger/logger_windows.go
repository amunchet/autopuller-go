//go:build windows
// +build windows

package logger

import (
	"io"
	"log"
	"os"
)

// InitLogger initializes the logger. It optionally logs to a specified file if a path is provided.
func InitLogger(logFilePath string) error {
	// Create a slice of writers; start with stdout
	writers := []io.Writer{os.Stdout}

	// If a log file path is provided, try to open or create the file
	if logFilePath != "" {
		logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		writers = append(writers, logFile)
	}

	// Combine all writers into one
	multiWriter := io.MultiWriter(writers...)

	// Set the log output to multiWriter (stdout + optional file)
	log.SetOutput(multiWriter)

	// Set log flags (date, time, file, line number)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Log an initialization message
	log.Println("Logger initialized (Windows)")

	return nil
}
