package logger

import (
	"bytes"
	"log"
	"testing"
)

func TestInitLogger(t *testing.T) {
	var buf bytes.Buffer

	// Set log output to buffer for testing
	logger := log.New(&buf, "", 0)
	logger.SetOutput(&buf)
	InitLogger()

	logger.Println("Test log message")
	logger.Println(buf.String())
	if !bytes.Contains(buf.Bytes(), []byte("Test log message")) {
		t.Errorf("Expected 'Test log message' to be logged, but it wasn't")
	}
}
