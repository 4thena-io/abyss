package database

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type SQLiteProvider struct {
	dbFile string
}

func NewSQLiteProvider() *SQLiteProvider {
	dbFile := "abyss.db"

	return &SQLiteProvider{
		dbFile: dbFile,
	}
}

func (p *SQLiteProvider) Connect() (*gorm.DB, error) {
	return gorm.Open(p.GetDialector(), &gorm.Config{})
}

func (p *SQLiteProvider) GetDialector() gorm.Dialector {
	return sqlite.Open(p.dbFile)
}
