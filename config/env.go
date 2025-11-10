package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	// Load .env file if it exists (ignore error if it doesn't)
	_ = godotenv.Load()
	return nil
}

// GetEnv gets an environment variable or returns a default value
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvInt gets an environment variable as int or returns a default value
func GetEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetEnvInt64 gets an environment variable as int64 or returns a default value
func GetEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetEnvDuration gets an environment variable as duration or returns a default value
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}

// MustGetEnv gets an environment variable or panics if it doesn't exist
func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}

// GetEnvBool gets an environment variable as bool or returns a default value
func GetEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}
