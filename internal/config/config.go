package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName          string
	AppEnv           string
	AppPort          string
	AppURL           string
	DatabaseDriver   string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseSSLMode  string
	SessionSecret    string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:          getEnv("APP_NAME", "Go Inertia Starter Kit"),
		AppEnv:           getEnv("APP_ENV", "local"),
		AppPort:          getEnv("APP_PORT", "8090"),
		AppURL:           getEnv("APP_URL", "http://localhost:8090"),
		DatabaseDriver:   getEnv("DATABASE_DRIVER", "postgres"),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseName:     getEnv("DATABASE_NAME", "starter_kit"),
		DatabaseUser:     getEnv("DATABASE_USER", "postgres"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "postgres"),
		DatabaseSSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
		SessionSecret:    getEnv("SESSION_SECRET", "super-secret-session-key-32-chars-long"),
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DatabaseHost,
		c.DatabasePort,
		c.DatabaseUser,
		c.DatabasePassword,
		c.DatabaseName,
		c.DatabaseSSLMode,
	)
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
