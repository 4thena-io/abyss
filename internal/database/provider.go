package database

import "gorm.io/gorm"

type Provider interface {
	Connect() (*gorm.DB, error)
	GetDialector() gorm.Dialector
}
