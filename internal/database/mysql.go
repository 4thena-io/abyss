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

func NewMySQLProvider() *MySQLProvider {
	return &MySQLProvider{
		host:     config.Environment.DbHost,
		port:     config.Environment.DbPort,
		user:     config.Environment.DbUser,
		password: config.Environment.DbPassword,
		dbName:   config.Environment.DbName,
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
