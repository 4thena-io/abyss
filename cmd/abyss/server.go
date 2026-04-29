package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	http *http.Server
}

type Config struct {
	Host string
	Port string
}

func newServer(cfg Config, handler http.Handler) *Server {
	addr := cfg.Host + ":" + cfg.Port

	return &Server{
		http: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) start() error {
	slog.Info("server running", "addr", "http://"+s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
