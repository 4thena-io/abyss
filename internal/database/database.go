package database

import (
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/config"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	schema       = config.Environment.DbSchema
	host         = config.Environment.DbHost
	port         = config.Environment.DbPort
	user         = config.Environment.DbUser
	password     = config.Environment.DbPassword
	databaseName = config.Environment.DbName
)

var Connection = initDatabaseConnection()

func initDatabaseConnection() *gorm.DB {
	// Step 1: Raw connection without search_path
	rawDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, databaseName,
	)

	rawDB, err := gorm.Open(postgres.Open(rawDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDB, err := rawDB.DB()
	if err != nil {
		panic(err)
	}

	if err := RunMigrations(sqlDB); err != nil {
		panic(fmt.Sprintf("Migration failed: %v", err))
	}

	finalDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		host, port, user, password, databaseName, schema,
	)

	db, err := gorm.Open(postgres.Open(finalDSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
