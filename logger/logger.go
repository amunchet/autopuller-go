package logger

import (
	"log"
	"os"
)

// InitLogger sets up the logger for the application.
func InitLogger() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)
	log.Println("Logger initialized")
}
