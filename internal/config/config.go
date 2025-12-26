package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
	Logger  LoggerConfig
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type StorageConfig struct {
	Type string
}

type LoggerConfig struct {
	Type  string
	Level string
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  parseDuration(getEnv("READ_TIMEOUT", "10s")),
			WriteTimeout: parseDuration(getEnv("WRITE_TIMEOUT", "10s")),
			IdleTimeout:  parseDuration(getEnv("IDLE_TIMEOUT", "60s")),
		},
		Storage: StorageConfig{
			getEnv("STORAGE_TYPE", "mem"),
		},
		Logger: LoggerConfig{
			Type:  getEnv("LOGGER_TYPE", "std"),
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseDuration(value string) time.Duration {
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return 10 * time.Second
}
