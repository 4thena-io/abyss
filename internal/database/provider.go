package database

import "gorm.io/gorm"

type Provider interface {
	GetDialector() gorm.Dialector
}
