package database

import (
	"errors"
	"fmt"
	"net"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/logger"
	"gorm.io/gorm"
)

// IsTransient reports connectivity/startup failures that can recover without
// configuration changes. Authentication and migration failures must not loop.
func IsTransient(err error) bool {
	var migration *MigrationError
	if errors.As(err, &migration) {
		return false
	}
	var state interface{ SQLState() string }
	if errors.As(err, &state) {
		code := state.SQLState()
		return len(code) >= 2 && (code[:2] == "08" || code == "57P03" || code == "53300")
	}
	var network net.Error
	return errors.As(err, &network)
}

type MigrationError struct{ Err error }

func (e *MigrationError) Error() string { return "database migration failed" }
func (e *MigrationError) Unwrap() error { return e.Err }

func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	provider, err := NewDatabaseProvider(cfg)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(provider.GetDialector(), &gorm.Config{
		Logger: logger.NewGormLogger(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := RunMigrations(db); err != nil {
		if sqlDB, closeErr := db.DB(); closeErr == nil {
			_ = sqlDB.Close()
		}
		return nil, &MigrationError{Err: err}
	}

	return db, nil
}
