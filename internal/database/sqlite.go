package database

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type SQLiteProvider struct {
	path string
}

func NewSQLiteProvider() *SQLiteProvider {
	return &SQLiteProvider{path: "abyss.db"}
}

func (p *SQLiteProvider) Connect() (*gorm.DB, error) {
	return gorm.Open(p.GetDialector(), &gorm.Config{})
}

func (p *SQLiteProvider) GetDialector() gorm.Dialector {
	return sqlite.Open(p.path)
}

