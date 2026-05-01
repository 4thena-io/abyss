package database

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
)

func NewDatabaseProvider(cfg config.DatabaseConfig) (Provider, error) {
	switch cfg.Type {
	case "postgres":
		return NewPostgresProvider(cfg), nil
	case "mysql":
		return NewMySQLProvider(cfg), nil
	case "sqlite", "":
		return NewSQLiteProvider(cfg), nil
	default:
		return nil, fmt.Errorf("unknown database provider: %s", cfg.Type)
	}
}
