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
	"github.com/luxe/backend/internal/pkg/imagekit"
	"github.com/luxe/backend/internal/pkg/response"
	"github.com/rs/zerolog/log"
)

type UploadHandler struct {
	uploadDir   string
	maxFileSize int64
	imagekit    *imagekit.Client
}

func NewUploadHandler(uploadDir string, maxFileSize int64, ik *imagekit.Client) *UploadHandler {
	_ = os.MkdirAll(uploadDir, 0755)
	return &UploadHandler{uploadDir: uploadDir, maxFileSize: maxFileSize, imagekit: ik}
}

func allowedImageExt(ext string) bool {
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg", ".png", ".webp", ".gif":
		return true
	default:
		return false
	}
}

// Upload handles multipart file uploads (ImageKit when configured, else local disk).
func (h *UploadHandler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxFileSize)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "no file provided or file too large")
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedImageExt(ext) {
		response.BadRequest(c, "only image files (jpg, jpeg, png, webp, gif) are allowed")
		return
	}

	if h.imagekit != nil && h.imagekit.Enabled() {
		result, err := h.imagekit.Upload(file, header.Filename)
		if err != nil {
			log.Error().Err(err).Msg("imagekit upload failed")
			response.InternalError(c, "imagekit upload failed: "+err.Error())
			return
		}
		log.Info().Str("url", result.URL).Msg("image uploaded to imagekit")
		response.Success(c, "file uploaded to imagekit", gin.H{
			"url":      result.URL,
			"filename": result.Name,
			"size":     result.Size,
			"file_id":  result.FileID,
			"provider": "imagekit",
		})
		return
	}

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

	publicURL := fmt.Sprintf("/uploads/%s/%s/%s", year, month, filename)
	response.Success(c, "file uploaded", gin.H{
		"url":      publicURL,
		"filename": filename,
		"size":     header.Size,
		"provider": "local",
	})
}

// UploadMultiple handles multiple file uploads.
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
		if !allowedImageExt(ext) {
			continue
		}

		if h.imagekit != nil && h.imagekit.Enabled() {
			src, err := fileHeader.Open()
			if err != nil {
				continue
			}
			result, err := h.imagekit.Upload(src, fileHeader.Filename)
			src.Close()
			if err != nil {
				continue
			}
			urls = append(urls, result.URL)
			continue
		}

		filename := fmt.Sprintf("%s%s", uuid.NewString(), ext)
		dst := filepath.Join(subDir, filename)
		if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
			continue
		}
		urls = append(urls, fmt.Sprintf("/uploads/%s/%s/%s", year, month, filename))
	}

	if len(urls) == 0 {
		response.BadRequest(c, "no files uploaded")
		return
	}

	response.Success(c, "files uploaded", gin.H{"urls": urls})
}
