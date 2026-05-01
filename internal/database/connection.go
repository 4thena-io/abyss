package database

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"gorm.io/gorm"
)

func NewConnection(cfg config.DatabaseConfig) (*gorm.DB, error) {
	provider, err := NewDatabaseProvider(cfg)
	if err != nil {
		return nil, err
	}

	db, err := provider.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := RunMigrations(db); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return db, nil
}
