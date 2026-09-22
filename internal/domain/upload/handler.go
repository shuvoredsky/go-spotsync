package upload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"spotsync/internal/httpresponse"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	MaxFileSize    = 5 * 1024 * 1024 // 5 MB max image upload size
	UploadsDirName = "uploads"
)

var allowedMIMETypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadImage(c echo.Context) error {
	// 1. Extract file from multipart form (supports form field "image" or "file")
	file, err := c.FormFile("image")
	if err != nil {
		file, err = c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, httpresponse.NewError("File upload error", "No image file provided in form-data. Please provide an image using key 'image' or 'file'"))
		}
	}

	// 2. Validate file size against maximum limit (5MB)
	if file.Size > MaxFileSize {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("File size exceeds limit", fmt.Sprintf("Uploaded file size (%d bytes) exceeds the maximum allowed limit of 5MB", file.Size)))
	}

	// 3. Open the uploaded source stream
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to open uploaded file", err.Error()))
	}
	defer src.Close()

	// 4. Sniff MIME type using the initial 512 bytes (content inspection to prevent spoofing)
	buffer := make([]byte, 512)
	n, err := src.Read(buffer)
	if err != nil && err != io.EOF {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to read file header", err.Error()))
	}

	detectedMIME := http.DetectContentType(buffer[:n])
	detectedMIME = strings.TrimSpace(strings.Split(detectedMIME, ";")[0])

	ext, isAllowed := allowedMIMETypes[detectedMIME]
	if !isAllowed {
		return c.JSON(http.StatusBadRequest, httpresponse.NewError("Invalid file format", fmt.Sprintf("MIME type '%s' is not supported. Only JPEG, PNG, WebP, and GIF images are allowed", detectedMIME)))
	}

	// Rewind file pointer back to start after MIME inspection
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to reset file stream", err.Error()))
	}

	// 5. Ensure upload directory exists securely
	if err := os.MkdirAll(UploadsDirName, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to initialize upload directory", err.Error()))
	}

	// 6. Generate cryptographically safe random filename (prevents directory traversal & collisions)
	randBytes := make([]byte, 8)
	if _, err := rand.Read(randBytes); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to generate secure file identifier", err.Error()))
	}
	safeFilename := fmt.Sprintf("img_%d_%s%s", time.Now().UnixMilli(), hex.EncodeToString(randBytes), ext)
	dstPath := filepath.Join(UploadsDirName, safeFilename)

	// 7. Write file to destination securely
	dst, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to create destination file", err.Error()))
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.NewError("Failed to write file data", err.Error()))
	}

	fileURL := fmt.Sprintf("/%s/%s", UploadsDirName, safeFilename)

	return c.JSON(http.StatusCreated, httpresponse.NewSuccess("Image uploaded successfully", map[string]any{
		"url":          fileURL,
		"filename":     safeFilename,
		"size":         file.Size,
		"content_type": detectedMIME,
	}))
}
