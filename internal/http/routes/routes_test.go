package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"clawbot-trust-lab/internal/http/handlers"
	"clawbot-trust-lab/internal/version"
)

func TestNewRoutesServesHealthAndVersion(t *testing.T) {
	router := New(func(next http.Handler) http.Handler { return next }, Services{
		System:   handlers.NewSystemHandler(func(context.Context) error { return nil }, version.Current()),
		TrustLab: &handlers.TrustLabHandler{},
		Operator: &handlers.OperatorHandler{},
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected health 200, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/version", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected version 200, got %d", rec.Code)
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode version payload: %v", err)
	}
	if payload["version"] == "" {
		t.Fatalf("expected non-empty version body %s", rec.Body.String())
	}
}
