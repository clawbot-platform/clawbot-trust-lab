package middleware

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestLoggerLogsMethodPathAndStatus(t *testing.T) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, nil))

	handler := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/scenarios/execute", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	body := buf.String()
	if !strings.Contains(body, "request completed") ||
		!strings.Contains(body, "method=POST") ||
		!strings.Contains(body, "path=/api/v1/scenarios/execute") ||
		!strings.Contains(body, "status=202") {
		t.Fatalf("unexpected log output %q", body)
	}
}
