package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	JWT      JWTConfig
	App      AppConfig
}

type ServerConfig struct {
	Port    string
	Mode    string // "debug" | "release"
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type JWTConfig struct {
	Secret          string
	ExpiryHours     int
	RefreshSecret   string
	RefreshExpHours int
}

type AppConfig struct {
	Name        string
	FrontendURL string
	UploadDir   string
	MaxFileSize int64 // bytes
}

// Load reads .env file and returns populated Config
func Load() *Config {
	// Load .env file (ignore error in production where env vars are injected)
	if err := godotenv.Load(); err != nil {
		log.Println("[config] No .env file found, using environment variables")
	}

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	refreshExpiry, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_HOURS", "168"))
	maxFileSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "10485760"), 10, 64) // 10MB default

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGO_DB", "luxe_db"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "luxe-super-secret-jwt-key-change-in-production"),
			ExpiryHours:     jwtExpiry,
			RefreshSecret:   getEnv("JWT_REFRESH_SECRET", "luxe-refresh-secret-key-change-in-production"),
			RefreshExpHours: refreshExpiry,
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "LUXE"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
			UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
			MaxFileSize: maxFileSize,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
