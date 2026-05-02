package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func serveCmd(configPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start the Abyss server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Resolve once so the setup handler writes to the exact
			// path this process (and the re-exec'd one) will load.
			if *configPath == "" {
				*configPath = config.DefaultPath()
			}

			cfg, err := config.Load(*configPath)
			if err != nil {
				return err
			}

			if err := logger.Init(cfg.Logging); err != nil {
				return err
			}

			restartCh := make(chan struct{}, 1)
			srv := setup(cfg, *configPath, restartCh)

			go func() {
				if err := srv.start(); err != nil && err != http.ErrServerClosed {
					log.Fatal().Err(err).Msg("server error")
				}
			}()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-quit:
				log.Info().Msg("shutting down server")
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				return srv.shutdown(ctx)

			case <-restartCh:
				log.Info().Msg("restarting with new configuration")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := srv.shutdown(ctx); err != nil {
					log.Warn().Err(err).Msg("error during shutdown before restart")
				}
				exe, err := os.Executable()
				if err != nil {
					return err
				}
				return syscall.Exec(exe, os.Args, os.Environ())
			}
		},
	}
}
