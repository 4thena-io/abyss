package logger

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	gormlogger "gorm.io/gorm/logger"
)

type gormLogger struct {
	logger        zerolog.Logger
	slowThreshold time.Duration
	level         gormlogger.LogLevel
}

func NewGormLogger() gormlogger.Interface {
	return &gormLogger{
		logger:        log.Logger.With().Str("component", "db").Logger(),
		slowThreshold: 200 * time.Millisecond,
		level:         gormlogger.Warn,
	}
}

func (l *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	clone := *l
	clone.level = level
	return &clone
}

func (l *gormLogger) Info(_ context.Context, msg string, args ...interface{}) {
	if l.level >= gormlogger.Info {
		l.logger.Info().Msgf(msg, args...)
	}
}

func (l *gormLogger) Warn(_ context.Context, msg string, args ...interface{}) {
	if l.level >= gormlogger.Warn {
		l.logger.Warn().Msgf(msg, args...)
	}
}

func (l *gormLogger) Error(_ context.Context, msg string, args ...interface{}) {
	if l.level >= gormlogger.Error {
		l.logger.Error().Msgf(msg, args...)
	}
}

func (l *gormLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.level <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && !errors.Is(err, gormlogger.ErrRecordNotFound):
		l.logger.Error().Err(err).Str("sql", sql).Int64("rows", rows).Dur("duration", elapsed).Msg("query error")
	case elapsed > l.slowThreshold:
		l.logger.Warn().Str("sql", sql).Int64("rows", rows).Dur("duration", elapsed).Msg("slow query")
	default:
		l.logger.Debug().Str("sql", sql).Int64("rows", rows).Dur("duration", elapsed).Msg("query")
	}
}
