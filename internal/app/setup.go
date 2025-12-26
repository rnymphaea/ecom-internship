package app

import (
	"log"

	"ecom-internship/internal/config"
	"ecom-internship/internal/database"
	"ecom-internship/internal/database/mem"
	"ecom-internship/internal/logger"
	"ecom-internship/internal/logger/std"
	"ecom-internship/internal/server"
)

type App struct {
	Server   *server.Server
	Database database.Database
	Logger   logger.Logger
}

func setup(cfg config.Config) (*App, error) {
	log, err := initLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	db, err := initDatabase(cfg.Storage, log)
	if err != nil {
		return nil, err
	}

	srv, err := initServer(cfg.Server, log, db)
	if err != nil {
		return nil, err
	}

	return &App{
		Server:   srv,
		Database: db,
		Logger:   log,
	}, nil
}

func initLogger(cfg config.LoggerConfig) (logger.Logger, error) {
	switch cfg.Type {
	case "std":
		return std.New(cfg.Level), nil
	default:
		log.Printf("unknown logger type: %s, using std", cfg.Type)
		return std.New(cfg.Level), nil
	}
}

func initDatabase(cfg config.StorageConfig, log logger.Logger) (database.Database, error) {
	switch cfg.Type {
	case "mem":
		return mem.New(log.With("component", "database")), nil
	default:
		log.Warn("unknown storage type: %s, using in-memory", cfg.Type)
		return mem.New(log.With("component", "database")), nil
	}
}

func initServer(cfg config.ServerConfig, log logger.Logger, db database.Database) (*server.Server, error) {
	router := server.NewRouter(log.With("component", "server"), db)
	srv := server.New(&cfg, router, log)
	return srv, nil
}
