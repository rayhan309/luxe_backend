package imagekit

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const uploadURL = "https://upload.imagekit.io/api/v1/files/upload"

// Config holds ImageKit credentials.
type Config struct {
	PublicKey   string
	PrivateKey  string
	URLEndpoint string
	Folder      string
}

// Client uploads files to ImageKit CDN.
type Client struct {
	privateKey  string
	publicKey   string
	urlEndpoint string
	folder      string
}

// UploadResult is returned after a successful upload.
type UploadResult struct {
	URL    string `json:"url"`
	FileID string `json:"fileId"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
}

// New creates an ImageKit client. Returns nil if private key is empty.
func New(cfg Config) *Client {
	if cfg.PrivateKey == "" {
		return nil
	}
	folder := cfg.Folder
	if folder == "" {
		folder = "/luxe"
	}
	if !strings.HasPrefix(folder, "/") {
		folder = "/" + folder
	}
	return &Client{
		privateKey:  cfg.PrivateKey,
		publicKey:   cfg.PublicKey,
		urlEndpoint: strings.TrimSuffix(cfg.URLEndpoint, "/"),
		folder:      folder,
	}
}

// Enabled reports whether ImageKit uploads are configured.
func (c *Client) Enabled() bool {
	return c != nil && c.privateKey != ""
}

// Upload sends a file buffer to ImageKit and returns the CDN URL.
func (c *Client) Upload(file io.Reader, originalFilename string) (*UploadResult, error) {
	if !c.Enabled() {
		return nil, fmt.Errorf("imagekit not configured")
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(originalFilename))
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copy file: %w", err)
	}

	safeName := sanitizeFilename(originalFilename)
	_ = writer.WriteField("fileName", safeName)
	_ = writer.WriteField("folder", c.folder)
	if c.publicKey != "" {
		_ = writer.WriteField("publicKey", c.publicKey)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, uploadURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.privateKey+":")))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("imagekit upload failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var result UploadResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	if result.URL == "" {
		return nil, fmt.Errorf("imagekit returned empty url")
	}

	return &result, nil
}

func sanitizeFilename(name string) string {
	base := filepath.Base(name)
	base = strings.ReplaceAll(base, " ", "-")
	return base
}
