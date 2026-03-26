package app

import (
	"bytes"
	"context"
	"testing"
)

func TestNewLoggerHonorsDebugLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("debug", &buf)
	logger.Debug("debug-enabled")

	if !bytes.Contains(buf.Bytes(), []byte("debug-enabled")) {
		t.Fatalf("expected debug log output, got %s", buf.String())
	}
	if !logger.Handler().Enabled(context.Background(), -4) {
		t.Fatal("expected debug level to be enabled")
	}
}
