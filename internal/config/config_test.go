package config

import (
	"testing"
	"time"
)

func TestLoad_DefaultValues(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 10*time.Second {
		t.Errorf("Expected ReadTimeout 10s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 10*time.Second {
		t.Errorf("Expected WriteTimeout 10s, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.Server.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout 60s, got %v", cfg.Server.IdleTimeout)
	}
	if cfg.Storage.Type != "mem" {
		t.Errorf("Expected storage type 'mem', got %s", cfg.Storage.Type)
	}
	if cfg.Logger.Type != "std" {
		t.Errorf("Expected logger type 'std', got %s", cfg.Logger.Type)
	}
	if cfg.Logger.Level != "info" {
		t.Errorf("Expected log level 'info', got %s", cfg.Logger.Level)
	}
}

func TestLoad_WithEnvironment(t *testing.T) {
	t.Setenv("PORT", "9000")
	t.Setenv("READ_TIMEOUT", "5s")
	t.Setenv("WRITE_TIMEOUT", "15s")
	t.Setenv("IDLE_TIMEOUT", "120s")
	t.Setenv("STORAGE_TYPE", "test")
	t.Setenv("LOGGER_TYPE", "custom")
	t.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 5*time.Second {
		t.Errorf("Expected ReadTimeout 5s, got %v", cfg.Server.ReadTimeout)
	}
	if cfg.Server.WriteTimeout != 15*time.Second {
		t.Errorf("Expected WriteTimeout 15s, got %v", cfg.Server.WriteTimeout)
	}
	if cfg.Server.IdleTimeout != 120*time.Second {
		t.Errorf("Expected IdleTimeout 120s, got %v", cfg.Server.IdleTimeout)
	}
	if cfg.Storage.Type != "test" {
		t.Errorf("Expected storage type 'test', got %s", cfg.Storage.Type)
	}
	if cfg.Logger.Type != "custom" {
		t.Errorf("Expected logger type 'custom', got %s", cfg.Logger.Type)
	}
	if cfg.Logger.Level != "debug" {
		t.Errorf("Expected log level 'debug', got %s", cfg.Logger.Level)
	}
}

func TestLoad_InvalidReadTimeout(t *testing.T) {
	t.Setenv("READ_TIMEOUT", "invalid")

	_, err := Load()
	if err == nil {
		t.Error("Expected error for invalid READ_TIMEOUT")
	}
}

func TestLoad_InvalidWriteTimeout(t *testing.T) {
	t.Setenv("WRITE_TIMEOUT", "invalid")

	_, err := Load()
	if err == nil {
		t.Error("Expected error for invalid WRITE_TIMEOUT")
	}
}

func TestLoad_InvalidIdleTimeout(t *testing.T) {
	t.Setenv("IDLE_TIMEOUT", "invalid")

	_, err := Load()
	if err == nil {
		t.Error("Expected error for invalid IDLE_TIMEOUT")
	}
}

func TestLoadServerConfig_Valid(t *testing.T) {
	t.Setenv("PORT", "3000")
	t.Setenv("READ_TIMEOUT", "30s")
	t.Setenv("WRITE_TIMEOUT", "20s")
	t.Setenv("IDLE_TIMEOUT", "90s")

	cfg, err := loadServerConfig()
	if err != nil {
		t.Fatalf("loadServerConfig failed: %v", err)
	}

	if cfg.Port != "3000" {
		t.Errorf("Expected port 3000, got %s", cfg.Port)
	}
	if cfg.ReadTimeout != 30*time.Second {
		t.Errorf("Expected ReadTimeout 30s, got %v", cfg.ReadTimeout)
	}
	if cfg.WriteTimeout != 20*time.Second {
		t.Errorf("Expected WriteTimeout 20s, got %v", cfg.WriteTimeout)
	}
	if cfg.IdleTimeout != 90*time.Second {
		t.Errorf("Expected IdleTimeout 90s, got %v", cfg.IdleTimeout)
	}
}

func TestLoadServerConfig_InvalidDuration(t *testing.T) {
	t.Setenv("READ_TIMEOUT", "not-a-duration")

	_, err := loadServerConfig()
	if err == nil {
		t.Error("Expected error for invalid duration")
	}
}

func TestLoadStorageConfig_Default(t *testing.T) {
	t.Setenv("STORAGE_TYPE", "mem")

	cfg, err := loadStorageConfig()
	if err != nil {
		t.Fatalf("loadStorageConfig failed: %v", err)
	}

	if cfg.Type != "mem" {
		t.Errorf("Expected storage type 'mem', got %s", cfg.Type)
	}
}

func TestLoadStorageConfig_WithEnv(t *testing.T) {
	t.Setenv("STORAGE_TYPE", "postgres")

	cfg, err := loadStorageConfig()
	if err != nil {
		t.Fatalf("loadStorageConfig failed: %v", err)
	}

	if cfg.Type != "postgres" {
		t.Errorf("Expected storage type 'postgres', got %s", cfg.Type)
	}
}

func TestLoadLoggerConfig_Default(t *testing.T) {
	t.Setenv("LOGGER_TYPE", "std")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := loadLoggerConfig()
	if err != nil {
		t.Fatalf("loadLoggerConfig failed: %v", err)
	}

	if cfg.Type != "std" {
		t.Errorf("Expected logger type 'std', got %s", cfg.Type)
	}
	if cfg.Level != "info" {
		t.Errorf("Expected log level 'info', got %s", cfg.Level)
	}
}

func TestLoadLoggerConfig_WithEnv(t *testing.T) {
	t.Setenv("LOGGER_TYPE", "json")
	t.Setenv("LOG_LEVEL", "warn")

	cfg, err := loadLoggerConfig()
	if err != nil {
		t.Fatalf("loadLoggerConfig failed: %v", err)
	}

	if cfg.Type != "json" {
		t.Errorf("Expected logger type 'json', got %s", cfg.Type)
	}
	if cfg.Level != "warn" {
		t.Errorf("Expected log level 'warn', got %s", cfg.Level)
	}
}

func TestGetEnv(t *testing.T) {
	if val := getEnv("TEST_KEY", "default"); val != "default" {
		t.Errorf("Expected 'default', got %s", val)
	}

	t.Setenv("TEST_KEY", "value")
	if val := getEnv("TEST_KEY", "default"); val != "value" {
		t.Errorf("Expected 'value', got %s", val)
	}

	t.Setenv("TEST_KEY", "")
	if val := getEnv("TEST_KEY", "default"); val != "" {
		t.Errorf("Expected empty string, got %s", val)
	}
}

func TestLoad_PartialErrors(t *testing.T) {
	t.Setenv("READ_TIMEOUT", "invalid")
	t.Setenv("WRITE_TIMEOUT", "invalid")
	t.Setenv("IDLE_TIMEOUT", "invalid")

	_, err := Load()
	if err == nil {
		t.Error("Expected error for invalid durations")
	}
}

func TestConfig_StructsNotEmpty(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("READ_TIMEOUT", "10s")
	t.Setenv("WRITE_TIMEOUT", "10s")
	t.Setenv("IDLE_TIMEOUT", "60s")
	t.Setenv("STORAGE_TYPE", "mem")
	t.Setenv("LOGGER_TYPE", "std")
	t.Setenv("LOG_LEVEL", "info")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server == nil {
		t.Error("Server config should not be nil")
	}
	if cfg.Storage == nil {
		t.Error("Storage config should not be nil")
	}
	if cfg.Logger == nil {
		t.Error("Logger config should not be nil")
	}
}

func TestLoad_EmptyEnvVars(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("STORAGE_TYPE", "")
	t.Setenv("LOGGER_TYPE", "")
	t.Setenv("LOG_LEVEL", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.Server.Port != "" {
		t.Errorf("Expected empty port, got %s", cfg.Server.Port)
	}
	if cfg.Storage.Type != "" {
		t.Errorf("Expected empty storage type, got %s", cfg.Storage.Type)
	}
	if cfg.Logger.Type != "" {
		t.Errorf("Expected empty logger type, got %s", cfg.Logger.Type)
	}
	if cfg.Logger.Level != "" {
		t.Errorf("Expected empty log level, got %s", cfg.Logger.Level)
	}
}
