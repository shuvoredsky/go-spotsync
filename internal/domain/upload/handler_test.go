package upload

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
)

// Minimal 1x1 valid PNG bytes
var validPNGBlob = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

func createMultipartRequest(t *testing.T, fieldName, fileName string, content []byte) (*http.Request, string) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	if _, err := io.Copy(part, bytes.NewReader(content)); err != nil {
		t.Fatalf("failed to write content: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/upload", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	return req, writer.FormDataContentType()
}

func TestUploadImage_Success(t *testing.T) {
	e := echo.New()
	req, _ := createMultipartRequest(t, "image", "avatar.png", validPNGBlob)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewUploadHandler()
	err := handler.UploadImage(c)
	if err != nil {
		t.Fatalf("handler returned unexpected error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201 Created, got %d. Body: %s", rec.Code, rec.Body.String())
	}

	var res map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to parse response JSON: %v", err)
	}

	if res["success"] != true {
		t.Errorf("expected success: true, got %v", res["success"])
	}

	data, ok := res["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data map in response, got %v", res["data"])
	}

	url, ok := data["url"].(string)
	if !ok || url == "" {
		t.Errorf("expected non-empty url, got %v", data["url"])
	}

	filename, ok := data["filename"].(string)
	if !ok || filename == "" {
		t.Errorf("expected non-empty filename, got %v", data["filename"])
	}

	// Clean up created test file
	uploadedFilePath := filepath.Join(UploadsDirName, filename)
	defer os.Remove(uploadedFilePath)

	if _, err := os.Stat(uploadedFilePath); os.IsNotExist(err) {
		t.Errorf("uploaded file does not exist on disk at %s", uploadedFilePath)
	}
}

func TestUploadImage_RejectDisguisedExecutable(t *testing.T) {
	e := echo.New()
	// An executable/script disguised as a jpg filename
	fakeImageBytes := []byte("#!/bin/bash\necho 'malicious script'\n")
	req, _ := createMultipartRequest(t, "image", "malicious.jpg", fakeImageBytes)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewUploadHandler()
	_ = handler.UploadImage(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request for invalid MIME, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestUploadImage_RejectExceedingSize(t *testing.T) {
	e := echo.New()
	// Buffer larger than 5MB
	largeBlob := make([]byte, 6*1024*1024)
	copy(largeBlob, validPNGBlob)

	req, _ := createMultipartRequest(t, "image", "huge.png", largeBlob)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewUploadHandler()
	_ = handler.UploadImage(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request for oversized file, got %d. Body: %s", rec.Code, rec.Body.String())
	}
}

func TestUploadImage_NoFileField(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/upload", bytes.NewReader([]byte("{}")))
	req.Header.Set(echo.HeaderContentType, "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := NewUploadHandler()
	_ = handler.UploadImage(c)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 Bad Request for missing form-data, got %d", rec.Code)
	}
}
