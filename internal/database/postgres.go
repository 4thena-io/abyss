package database

import (
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresProvider struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func NewPostgresProvider() *PostgresProvider {
	return &PostgresProvider{
		host:     config.Environment.DbHost,
		port:     config.Environment.DbPort,
		user:     config.Environment.DbUser,
		password: config.Environment.DbPassword,
		dbName:   config.Environment.DbName,
	}
}

func (p *PostgresProvider) Connect() (*gorm.DB, error) {
	return gorm.Open(p.GetDialector(), &gorm.Config{})
}

func (p *PostgresProvider) GetDialector() gorm.Dialector {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		p.host, p.port, p.user, p.password, p.dbName,
	)

	return postgres.Open(dsn)
}
