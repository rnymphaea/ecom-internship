// Package server provides HTTP server implementation.
package server

import (
	"context"
	"net/http"

	"ecom-internship/internal/config"
	"ecom-internship/internal/logger"
)

// Server represents the HTTP server.
type Server struct {
	server *http.Server
	log    logger.Logger
}

// New creates a new HTTP server instance.
func New(cfg *config.ServerConfig, router *http.ServeMux, log logger.Logger) *Server {
	return &Server{
		server: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        router,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			IdleTimeout:    cfg.IdleTimeout,
			MaxHeaderBytes: 1 << 20,
		},
		log: log,
	}
}

// Start begins listening for HTTP requests.
func (s *Server) Start() error {
	s.log.Info("starting HTTP server", "port", s.server.Addr)

	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server.
func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("shutting down server")

	return s.server.Shutdown(ctx)
}
