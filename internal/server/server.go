package server

import (
	"cakes-database-app/internal/config"
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	HTTPServer *http.Server
}

func (s *Server) Run(handler http.Handler, cfg *config.Config) error {

	timeout, err := time.ParseDuration(cfg.HTTPServer.Timeout)
	if err != nil {
		log.Fatalf("invalid timeout configuration: %s", err)
	}
	idleTimeout, err := time.ParseDuration(cfg.HTTPServer.IdleTimeout)
	if err != nil {
		log.Fatalf("invalid idle timeout configuration: %s", err)
	}

	s.HTTPServer = &http.Server{
		Handler:     handler,
		Addr:        cfg.HTTPServer.Address,
		ReadTimeout: timeout,
		IdleTimeout: idleTimeout,
	}

	return s.HTTPServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.HTTPServer.Shutdown(ctx)
}
