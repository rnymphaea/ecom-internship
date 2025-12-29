// Package app provides the main application initialization and shutdown logic.
package app

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ecom-internship/internal/config"
)

// Run initializes and starts the HTTP server with graceful shutdown.
func Run() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	app, err := setup(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app.Logger.Info("config loaded successfully")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.Logger.Error("failed to start server", "error", err)
		}
	}()

	<-done
	app.Logger.Info("server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.Server.Stop(ctx); err != nil {
		app.Logger.Error("failed to shutdown server", "error", err)
	}

	app.Logger.Info("server stopped gracefully")
}
