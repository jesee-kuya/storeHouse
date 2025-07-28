package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	// Server Configuration
	Server ServerConfig

	// Database Configuration
	Database DatabaseConfig

	// SMS Configuration
	SMS SMSConfig

	// JWT Configuration
	JWT JWTConfig

	// App Configuration
	App AppConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	Host         string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Path            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	MigrationsPath  string
}

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	Provider    string // "africastalking" or "twilio"
	APIKey      string
	APISecret   string
	Username    string
	SenderID    string
	Environment string // "sandbox" or "production"
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	SecretKey      string
	TokenDuration  time.Duration
	RefreshDuration time.Duration
	Issuer         string
}

// AppConfig holds general application configuration
type AppConfig struct {
	Name        string
	Version     string
	Debug       bool
	LogLevel    string
	ChurchName  string
	ChurchPhone string
	Currency    string
	Timezone    string
}

// Global config instance
//var AppConfig *Config

// Load loads configuration from environment variables and .env file
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "localhost"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			ReadTimeout:  parseDuration(getEnv("SERVER_READ_TIMEOUT", "15s")),
			WriteTimeout: parseDuration(getEnv("SERVER_WRITE_TIMEOUT", "15s")),
			IdleTimeout:  parseDuration(getEnv("SERVER_IDLE_TIMEOUT", "60s")),
		},
		Database: DatabaseConfig{
			Path:            getEnv("DB_PATH", "./data/church_app.db"),
			MaxOpenConns:    parseInt(getEnv("DB_MAX_OPEN_CONNS", "25")),
			MaxIdleConns:    parseInt(getEnv("DB_MAX_IDLE_CONNS", "25")),
			ConnMaxLifetime: parseDuration(getEnv("DB_CONN_MAX_LIFETIME", "5m")),
			MigrationsPath:  getEnv("DB_MIGRATIONS_PATH", "./database/migrations"),
		},
		SMS: SMSConfig{
			Provider:    getEnv("SMS_PROVIDER", "africastalking"),
			APIKey:      getEnv("SMS_API_KEY", ""),
			APISecret:   getEnv("SMS_API_SECRET", ""),
			Username:    getEnv("SMS_USERNAME", "sandbox"),
			SenderID:    getEnv("SMS_SENDER_ID", "CHURCH"),
			Environment: getEnv("SMS_ENVIRONMENT", "sandbox"),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET_KEY", generateDefaultJWTSecret()),
			TokenDuration:   parseDuration(getEnv("JWT_TOKEN_DURATION", "24h")),
			RefreshDuration: parseDuration(getEnv("JWT_REFRESH_DURATION", "168h")), // 7 days
			Issuer:          getEnv("JWT_ISSUER", "church-receipts-app"),
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "Church Receipts & Finance"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Debug:       parseBool(getEnv("DEBUG", "true")),
			LogLevel:    getEnv("LOG_LEVEL", "info"),
			ChurchName:  getEnv("CHURCH_NAME", "Sample Church"),
			ChurchPhone: getEnv("CHURCH_PHONE", "+254700000000"),
			Currency:    getEnv("CURRENCY", "KES"),
			Timezone:    getEnv("TIMEZONE", "Africa/Nairobi"),
		},
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Set global config
	// AppConfig = config

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server configuration
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate JWT secret in production
	if c.Server.Environment == "production" && c.JWT.SecretKey == generateDefaultJWTSecret() {
		return fmt.Errorf("JWT_SECRET_KEY must be set in production environment")
	}

	// Validate SMS configuration if provider is set
	if c.SMS.Provider != "" {
		switch c.SMS.Provider {
		case "africastalking":
			if c.SMS.APIKey == "" {
				return fmt.Errorf("SMS_API_KEY is required for Africa's Talking")
			}
			if c.SMS.Username == "" {
				return fmt.Errorf("SMS_USERNAME is required for Africa's Talking")
			}
		case "twilio":
			if c.SMS.APIKey == "" {
				return fmt.Errorf("SMS_API_KEY (Account SID) is required for Twilio")
			}
			if c.SMS.APISecret == "" {
				return fmt.Errorf("SMS_API_SECRET (Auth Token) is required for Twilio")
			}
		default:
			return fmt.Errorf("unsupported SMS provider: %s", c.SMS.Provider)
		}
	}

	// Validate database path directory exists or can be created
	if err := ensureDirectoryExists(c.Database.Path); err != nil {
		return fmt.Errorf("database path validation failed: %w", err)
	}

	return nil
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return c.Server.Host + ":" + c.Server.Port
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	return c.Database.Path
}

// PrintConfig prints the current configuration (safe for logging)
func (c *Config) PrintConfig() {
	log.Println("=== Application Configuration ===")
	log.Printf("App Name: %s", c.App.Name)
	log.Printf("Version: %s", c.App.Version)
	log.Printf("Environment: %s", c.Server.Environment)
	log.Printf("Server Address: %s", c.GetServerAddress())
	log.Printf("Database Path: %s", c.Database.Path)
	log.Printf("SMS Provider: %s", c.SMS.Provider)
	log.Printf("Church Name: %s", c.App.ChurchName)
	log.Printf("Currency: %s", c.App.Currency)
	log.Printf("Timezone: %s", c.App.Timezone)
	log.Printf("Debug Mode: %t", c.App.Debug)
	log.Println("================================")
}

// Helper functions

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// parseInt parses a string to int with error handling
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("Warning: failed to parse int '%s', using 0", s)
		return 0
	}
	return i
}

// parseBool parses a string to bool with error handling
func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Printf("Warning: failed to parse bool '%s', using false", s)
		return false
	}
	return b
}

// parseDuration parses a string to time.Duration with error handling
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("Warning: failed to parse duration '%s', using 0", s)
		return 0
	}
	return d
}

// generateDefaultJWTSecret generates a default JWT secret (should not be used in production)
func generateDefaultJWTSecret() string {
	return "church-receipts-default-secret-change-in-production"
}

// ensureDirectoryExists ensures the directory for the given file path exists
func ensureDirectoryExists(filePath string) error {
	dir := filepath.Dir(filePath)
	return os.MkdirAll(dir, 0755)
}

// GetConfig returns the global configuration instance
func GetConfig() *Config {
	// if AppConfig == nil {
	// 	log.Fatal("Configuration not loaded. Call config.Load() first.")
	// }
	// return AppConfig
	return &Config{}
}

// Example .env file content for reference
const ExampleEnvFile = `
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost
ENVIRONMENT=development

# Database Configuration
DB_PATH=./data/church_app.db
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25

# SMS Configuration (Africa's Talking)
SMS_PROVIDER=africastalking
SMS_API_KEY=your_api_key_here
SMS_USERNAME=your_username_here
SMS_SENDER_ID=CHURCH
SMS_ENVIRONMENT=sandbox

# SMS Configuration (Twilio Alternative)
# SMS_PROVIDER=twilio
# SMS_API_KEY=your_account_sid
# SMS_API_SECRET=your_auth_token
# SMS_SENDER_ID=+1234567890

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key-here
JWT_TOKEN_DURATION=24h
JWT_REFRESH_DURATION=168h

# Application Configuration
APP_NAME=Church Receipts & Finance
APP_VERSION=1.0.0
DEBUG=true
LOG_LEVEL=info
CHURCH_NAME=Your Church Name
CHURCH_PHONE=+254700000000
CURRENCY=KES
TIMEZONE=Africa/Nairobi
`