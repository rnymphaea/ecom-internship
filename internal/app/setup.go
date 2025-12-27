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

func setup(cfg *config.Config) (*App, error) {
	rootLogger, err := initLogger(cfg.Logger)
	if err != nil {
		return nil, err
	}

	dbLogger := rootLogger.With("component", "database")
	db, err := initDatabase(cfg.Storage, dbLogger)
	if err != nil {
		return nil, err
	}

	srvLogger := rootLogger.With("component", "database")
	srv, err := initServer(cfg.Server, srvLogger, db)
	if err != nil {
		return nil, err
	}

	return &App{
		Server:   srv,
		Database: db,
		Logger:   rootLogger,
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
		return mem.New(log), nil
	default:
		log.Warn("unknown storage type, using in-memory", "type", cfg.Type)
		return mem.New(log), nil
	}
}

func initServer(cfg config.ServerConfig, log logger.Logger, db database.Database) (*server.Server, error) {
	router := server.NewRouter(log, db)
	srv := server.New(&cfg, router, log)
	return srv, nil
}
