package logger

import (
	"io"
	"os"
	"time"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(cfg config.LoggingConfig) error {
	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	var stdoutWriter io.Writer
	if cfg.Format == "pretty" {
		stdoutWriter = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	} else {
		stdoutWriter = os.Stdout
	}

	if cfg.File != "" {
		f, err := os.OpenFile(cfg.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		log.Logger = zerolog.New(io.MultiWriter(stdoutWriter, f)).
			With().
			Timestamp().
			Str("app", "abyss").
			Logger()
		return nil
	}

	log.Logger = zerolog.New(stdoutWriter).
		With().
		Timestamp().
		Str("app", "abyss").
		Logger()

	return nil
}
