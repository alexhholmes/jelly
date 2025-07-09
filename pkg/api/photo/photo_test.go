package photo

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"jelly/pkg/api/gen"
	"jelly/pkg/api/util"
)

func TestPhotoHandler_UploadPhoto_Success(t *testing.T) {
	// Set test environment variable for max file size
	originalEnv := os.Getenv("PHOTO_MAX_FILE_SIZE_MB")
	os.Setenv("PHOTO_MAX_FILE_SIZE_MB", "10")
	defer func() {
		if originalEnv == "" {
			os.Unsetenv("PHOTO_MAX_FILE_SIZE_MB")
		} else {
			os.Setenv("PHOTO_MAX_FILE_SIZE_MB", originalEnv)
		}
	}()
	
	handler := PhotoHandler{}

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	fileWriter, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	fileWriter.Write([]byte("fake image data"))

	// Add caption
	writer.WriteField("caption", "Test caption")

	// Add tags
	writer.WriteField("tags", "tag1")
	writer.WriteField("tags", "tag2")

	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add logger to context
	logger := slog.Default()
	ctx := context.WithValue(req.Context(), util.ContextLogger, logger)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.UploadPhoto(w, req)

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
	var resp gen.PhotoUploadResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify photo response structure
	if resp.Photo.Id == "" {
		t.Error("Expected photo ID to be set")
	}

	if resp.Photo.Url == "" {
		t.Error("Expected photo URL to be set")
	}

	if resp.Photo.Caption == nil || *resp.Photo.Caption != "Test caption" {
		t.Errorf("Expected caption 'Test caption', got %v", resp.Photo.Caption)
	}

	if resp.Photo.Tags == nil || len(*resp.Photo.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %v", resp.Photo.Tags)
	} else {
		tags := *resp.Photo.Tags
		if tags[0] != "tag1" || tags[1] != "tag2" {
			t.Errorf("Expected tags [tag1, tag2], got %v", tags)
		}
	}

	if resp.Photo.UploadedAt.IsZero() {
		t.Error("Expected uploadedAt to be set")
	}

	if resp.Message == nil || *resp.Message != "Photo uploaded successfully" {
		t.Errorf("Expected message 'Photo uploaded successfully', got %v", resp.Message)
	}
}

func TestPhotoHandler_UploadPhoto_NoFile(t *testing.T) {
	handler := PhotoHandler{}

	// Create multipart form data without file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("caption", "Test caption")
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add logger to context
	logger := slog.Default()
	ctx := context.WithValue(req.Context(), util.ContextLogger, logger)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.UploadPhoto(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check error message
	if !strings.Contains(w.Body.String(), util.ErrMsgFileRequired) {
		t.Errorf("Expected error message about file being required, got %s", w.Body.String())
	}
}

func TestPhotoHandler_UploadPhoto_InvalidForm(t *testing.T) {
	handler := PhotoHandler{}

	// Create invalid form data
	body := strings.NewReader("invalid form data")
	req := httptest.NewRequest(http.MethodPost, "/photo", body)
	req.Header.Set("Content-Type", "multipart/form-data; boundary=invalid")

	// Add logger to context
	logger := slog.Default()
	ctx := context.WithValue(req.Context(), util.ContextLogger, logger)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.UploadPhoto(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Check error message
	if !strings.Contains(w.Body.String(), util.ErrMsgFailedToParseForm) {
		t.Errorf("Expected error message about parsing form, got %s", w.Body.String())
	}
}

func TestPhotoHandler_UploadPhoto_WithoutLogger(t *testing.T) {
	handler := PhotoHandler{}

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, _ := writer.CreateFormFile("file", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()

	// This should panic due to missing logger in context
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when logger is missing from context")
		}
	}()

	handler.UploadPhoto(w, req)
}

func TestPhotoHandler_UploadPhoto_MinimalData(t *testing.T) {
	handler := PhotoHandler{}

	// Create multipart form data with only file
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, err := writer.CreateFormFile("file", "minimal.jpg")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	fileWriter.Write([]byte("minimal image data"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/photo", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Add logger to context
	logger := slog.Default()
	ctx := context.WithValue(req.Context(), util.ContextLogger, logger)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.UploadPhoto(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var resp gen.PhotoUploadResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify minimal response - caption should be empty string, tags should be empty array
	if resp.Photo.Caption == nil || *resp.Photo.Caption != "" {
		t.Errorf("Expected empty caption, got %v", resp.Photo.Caption)
	}

	if resp.Photo.Tags == nil || len(*resp.Photo.Tags) != 0 {
		t.Errorf("Expected empty tags array, got %v", resp.Photo.Tags)
	}
}
