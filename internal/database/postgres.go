package database

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
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

func NewPostgresProvider(cfg config.DatabaseConfig) *PostgresProvider {
	return &PostgresProvider{
		host:     cfg.Host,
		port:     cfg.Port,
		user:     cfg.User,
		password: cfg.Password,
		dbName:   cfg.Name,
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
