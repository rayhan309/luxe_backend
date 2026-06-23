package http

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luxe/backend/internal/pkg/response"
)

type UploadHandler struct {
	uploadDir   string
	maxFileSize int64
}

func NewUploadHandler(uploadDir string, maxFileSize int64) *UploadHandler {
	// Ensure upload directory exists
	_ = os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{uploadDir: uploadDir, maxFileSize: maxFileSize}
}

// Upload handles multipart file uploads
func (h *UploadHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxFileSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "no file provided or file too large")
		return
	}
	defer file.Close()

	// Validate extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowedExts[ext] {
		response.BadRequest(c, "only image files (jpg, jpeg, png, webp, gif) are allowed")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s-%d%s", uuid.NewString(), time.Now().UnixNano(), ext)
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	subDir := filepath.Join(h.uploadDir, year, month)
	_ = os.MkdirAll(subDir, 0755)

	dst := filepath.Join(subDir, filename)
	if err := c.SaveUploadedFile(header, dst); err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	// Return public URL (relative path)
	publicURL := fmt.Sprintf("/uploads/%s/%s/%s", year, month, filename)
	response.Success(c, "file uploaded", gin.H{
		"url":      publicURL,
		"filename": filename,
		"size":     header.Size,
	})
}

// UploadMultiple handles multiple file uploads
func (h *UploadHandler) UploadMultiple(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "invalid multipart form")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.BadRequest(c, "no files provided")
		return
	}

	var urls []string
	year := time.Now().Format("2006")
	month := time.Now().Format("01")
	subDir := filepath.Join(h.uploadDir, year, month)
	_ = os.MkdirAll(subDir, 0755)

	for _, fileHeader := range files {
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
		if !allowedExts[ext] {
			continue
		}

		filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
		dst := filepath.Join(subDir, filename)
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			continue
		}
		urls = append(urls, fmt.Sprintf("/uploads/%s/%s/%s", year, month, filename))
	}

	response.Success(c, "files uploaded", gin.H{"urls": urls})
}
