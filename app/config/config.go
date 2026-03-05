package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	JWT            JWTConfig
	Supabase       SupabaseConfig
	Payment        PaymentConfig
	Email          EmailConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

type SupabaseConfig struct {
	URL    string
	Key    string
}

type PaymentConfig struct {
	XenditSecretKey string
	XenditPublicKey string
	IsProduction      bool
}

type EmailConfig struct {
	GmailAPIKey            string
	GmailServiceAccountJSON string
	FromEmail              string
	FromName               string
	IsEnabled              bool
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "rental_car"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			Expiration: getEnvInt("JWT_EXPIRATION", 24),
		},
		Supabase: SupabaseConfig{
			URL:    getEnv("SUPABASE_URL", ""),
			Key:    getEnv("SUPABASE_KEY", ""),
		},
		Payment: PaymentConfig{
			XenditSecretKey:   getEnv("XENDIT_SECRET_KEY", ""),
			XenditPublicKey:   getEnv("XENDIT_PUBLIC_KEY", ""),
			IsProduction:      getEnv("ENV", "development") == "production",
		},
		Email: EmailConfig{
			GmailAPIKey:             getEnv("GMAIL_API_KEY", ""),
			GmailServiceAccountJSON: getEnv("GMAIL_SERVICE_ACCOUNT_JSON", ""),
			FromEmail:               getEnv("SMTP_FROM_EMAIL", "noreply@rentalcar.com"),
			FromName:                getEnv("EMAIL_FROM_NAME", "Rental Car Service"),
			IsEnabled:               getEnv("EMAIL_ENABLED", "false") == "true",
		},
	}
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
