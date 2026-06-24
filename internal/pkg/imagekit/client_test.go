package imagekit

import (
	"bytes"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestUploadIntegration(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	privateKey := os.Getenv("IMAGEKIT_PRIVATE_KEY")
	if privateKey == "" {
		t.Skip("IMAGEKIT_PRIVATE_KEY not set")
	}

	client := New(Config{
		PublicKey:   os.Getenv("IMAGEKIT_PUBLIC_KEY"),
		PrivateKey:  privateKey,
		URLEndpoint: os.Getenv("IMAGEKIT_URL_ENDPOINT"),
		Folder:      "/luxe/test",
	})

	// 1x1 PNG
	png := []byte{
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

	result, err := client.Upload(bytes.NewReader(png), "test-upload.png")
	if err != nil {
		t.Fatalf("upload failed: %v", err)
	}
	if result.URL == "" || !contains(result.URL, "ik.imagekit.io") {
		t.Fatalf("expected imagekit url, got: %s", result.URL)
	}
	t.Logf("uploaded: %s", result.URL)
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
