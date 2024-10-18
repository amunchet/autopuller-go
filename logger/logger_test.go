package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestInitLogger(t *testing.T) {
	var buf bytes.Buffer

	// Set log output to buffer for testing
	log.SetOutput(&buf)
	InitLogger()

	log.Println("Test log message")

	if !bytes.Contains(buf.Bytes(), []byte("Test log message")) {
		t.Errorf("Expected 'Test log message' to be logged, but it wasn't")
	}
}
