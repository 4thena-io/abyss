package database

import (
	"github.com/4thena-io/abyss/internal/config"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SQLiteProvider struct {
	path string
}

func NewSQLiteProvider(cfg config.DatabaseConfig) *SQLiteProvider {
	path := cfg.Path
	if path == "" {
		path = "abyss.db"
	}
	return &SQLiteProvider{path: path}
}

func (p *SQLiteProvider) Connect() (*gorm.DB, error) {
	return gorm.Open(p.GetDialector(), &gorm.Config{})
}

func (p *SQLiteProvider) GetDialector() gorm.Dialector {
	return sqlite.Open(p.path)
}
