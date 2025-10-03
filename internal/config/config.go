package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	Auth          AuthConfig
	Upload        UploadConfig
	Notifications NotificationConfig
	Frontend      FrontendConfig
}

type ServerConfig struct {
	Host string
	Port int
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}


type AuthConfig struct {
	JWTSecret string
}

type UploadConfig struct {
	MaxFileSize int64  // in bytes
	UploadDir   string
	AllowedTypes []string
}

type NotificationConfig struct {
	EmailServiceURL   string
	EmailServiceToken string
	SlackServiceURL   string
	SlackServiceToken string
	AdminEmails       []string
	AdminEmail        string
	SlackChannel      string
}

type FrontendConfig struct {
	BaseURL string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "localhost"),
			Port: getEnvAsInt("SERVER_PORT", 8080),
			Env:  getEnv("SERVER_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "community_support"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Auth: AuthConfig{
			JWTSecret: getEnv("JWT_SECRET", "your-secret-key"),
		},
		Notifications: NotificationConfig{
			EmailServiceURL:   getEnv("EMAIL_SERVICE_URL", "http://localhost:8081"),
			EmailServiceToken: getEnv("EMAIL_SERVICE_TOKEN", ""),
			SlackServiceURL:   getEnv("SLACK_SERVICE_URL", "http://localhost:8082"),
			SlackServiceToken: getEnv("SLACK_SERVICE_TOKEN", ""),
			AdminEmails:       getEnvAsStringSlice("ADMIN_EMAILS", []string{"admin@kemuko.com"}),
			AdminEmail:        getEnv("ADMIN_EMAIL", "support@kemuko.com"),
			SlackChannel:      getEnv("SLACK_CHANNEL", "#kemuko-support"),
		},
		Frontend: FrontendConfig{
			BaseURL: getEnv("FRONTEND_BASE_URL", "https://kemuko.com"),
		},
		Upload: UploadConfig{
			MaxFileSize: getEnvAsInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB
			UploadDir:   getEnv("UPLOAD_DIR", "./uploads"),
			AllowedTypes: []string{
				"image/jpeg", "image/png", "image/gif",
				"application/pdf", "text/plain",
				"application/zip", "application/x-zip-compressed",
			},
		},
	}

	return config, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Split by comma and trim spaces
		parts := []string{}
		for _, part := range strings.Split(value, ",") {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				parts = append(parts, trimmed)
			}
		}
		if len(parts) > 0 {
			return parts
		}
	}
	return defaultValue
}