package database

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"gorm.io/gorm"
)

var Connection = initDatabaseConnection()

func initDatabaseConnection() *gorm.DB {
	provider, err := NewDatabaseProvider(config.Environment.DbType)
	if err != nil {
		panic(err)
	}

	db, err := provider.Connect()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	if err := RunMigrations(db); err != nil {
		panic(fmt.Sprintf("migration failed: %v", err))
	}

	return db
}
