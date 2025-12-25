package config

import (
	"os"
	"time"
)

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
