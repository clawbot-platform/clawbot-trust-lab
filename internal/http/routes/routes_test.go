package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
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
	if !strings.Contains(rec.Body.String(), "dev") {
		t.Fatalf("unexpected version body %s", rec.Body.String())
	}
}
