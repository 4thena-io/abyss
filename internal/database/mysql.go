package database

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLProvider struct {
	host     string
	port     string
	user     string
	password string
	dbName   string
}

func NewMySQLProvider(cfg config.DatabaseConfig) *MySQLProvider {
	return &MySQLProvider{
		host:     cfg.Host,
		port:     cfg.Port,
		user:     cfg.User,
		password: cfg.Password,
		dbName:   cfg.Name,
	}
}

func (p *MySQLProvider) Connect() (*gorm.DB, error) {
	return gorm.Open(p.GetDialector(), &gorm.Config{})
}

func (p *MySQLProvider) GetDialector() gorm.Dialector {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		p.user, p.password, p.host, p.port, p.dbName,
	)
	return mysql.Open(dsn)
}
