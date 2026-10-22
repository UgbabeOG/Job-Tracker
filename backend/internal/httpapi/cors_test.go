package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddlewareAllowsConfiguredOrigin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORSMiddleware(next, []string{
		"http://localhost:5173",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("expected allowed origin header, got %q", got)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestCORSMiddlewareDoesNotAllowUnconfiguredOrigin(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := CORSMiddleware(next, []string{
		"http://localhost:5173",
	})

	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	req.Header.Set("Origin", "https://untrusted.example")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no allow-origin header, got %q", got)
	}
}

func TestCORSMiddlewareHandlesPreflight(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("preflight request should not reach the next handler")
	})

	handler := CORSMiddleware(next, []string{
		"http://localhost:5173",
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/jobs", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "PATCH")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("unexpected allow-origin header: %q", got)
	}

	if got := rec.Header().Get("Access-Control-Allow-Methods"); got == "" {
		t.Fatal("expected allowed methods header")
	}

	if got := rec.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Fatal("expected allowed headers")
	}
}

func TestCORSMiddlewareRejectsUnsupportedPreflightMethod(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("preflight request should not reach the next handler")
	})

	handler := CORSMiddleware(next, []string{
		"http://localhost:5173",
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/jobs", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "PUT")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", rec.Code)
	}
}
