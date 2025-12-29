// Package config provides application configuration loading.
package config

import (
	"os"
	"time"
)

// Config contains all application configuration.
type Config struct {
	Server  *ServerConfig
	Storage *StorageConfig
	Logger  *LoggerConfig
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// StorageConfig contains data storage settings.
type StorageConfig struct {
	Type string
}

// LoggerConfig contains logger settings.
type LoggerConfig struct {
	Type  string
	Level string
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	server, err := loadServerConfig()
	if err != nil {
		return nil, err
	}

	storage, err := loadStorageConfig()
	if err != nil {
		return nil, err
	}

	logger, err := loadLoggerConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server:  server,
		Storage: storage,
		Logger:  logger,
	}

	return cfg, nil
}

func loadServerConfig() (*ServerConfig, error) {
	port := getEnv("PORT", "8080")

	readTimeout, err := time.ParseDuration(getEnv("READ_TIMEOUT", "10s"))
	if err != nil {
		return nil, err
	}

	writeTimeout, err := time.ParseDuration(getEnv("WRITE_TIMEOUT", "10s"))
	if err != nil {
		return nil, err
	}

	idleTimeout, err := time.ParseDuration(getEnv("IDLE_TIMEOUT", "60s"))
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}, nil
}

//nolint:unparam
func loadStorageConfig() (*StorageConfig, error) {
	return &StorageConfig{
		Type: getEnv("STORAGE_TYPE", "mem"),
	}, nil
}

//nolint:unparam
func loadLoggerConfig() (*LoggerConfig, error) {
	return &LoggerConfig{
		Type:  getEnv("LOGGER_TYPE", "std"),
		Level: getEnv("LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
