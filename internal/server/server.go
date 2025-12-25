package server

import (
	"context"
	"net/http"

	"ecom-internship/internal/config"
	"ecom-internship/internal/logger"
)

type Server struct {
	server *http.Server
	log    logger.Logger
}

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

func (s *Server) Start() error {
	s.log.Debugf("starting HTTP server")
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("shutting down server")
	return s.server.Shutdown(ctx)
}
