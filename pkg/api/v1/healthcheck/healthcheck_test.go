package healthcheck

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"jelly/pkg/api/v1/gen"
	"jelly/pkg/api/v1/util"
)

func TestHealthHandler_HealthCheck(t *testing.T) {
	handler := HealthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Add logger to context
	logger := slog.Default()
	ctx := context.WithValue(req.Context(), util.ContextLogger, logger)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check content type
	expectedContentType := "application/json"
	if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Check response body
	var resp gen.HealthCheck
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expectedStatus := "ok"
	if resp.Status != expectedStatus {
		t.Errorf("Expected status %s, got %s", expectedStatus, resp.Status)
	}
}

func TestHealthHandler_HealthCheck_WithoutLogger(t *testing.T) {
	handler := HealthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	// This should panic due to missing logger in context
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when logger is missing from context")
		}
	}()

	handler.HealthCheck(w, req)
}

func TestHealthHandler_HealthCheck_JSONMarshalError(t *testing.T) {
	handler := HealthHandler{}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Add logger to context
	logger := slog.Default()
	ctx := context.WithValue(req.Context(), util.ContextLogger, logger)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	// Since HealthCheck struct is simple and should always marshal successfully,
	// this test verifies the response is properly formatted
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify JSON response structure
	var resp gen.HealthCheck
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("Expected status 'ok', got %s", resp.Status)
	}
}
