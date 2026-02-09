package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/4thena-io/abyss/internal/api/rest/router"
)

type Server struct {
	http *http.Server
}

type Config struct {
	Host string
	Port string
}

func New(cfg Config) *Server {
	addr := cfg.Host + ":" + cfg.Port

	return &Server{
		http: &http.Server{
			Addr:         addr,
			Handler:      router.New(),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	slog.Info("server running", "addr", "http://"+s.http.Addr)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
