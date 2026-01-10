package database

import (
	"git.4thena.io/4thena/abys/internal/model"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Project{},
		&model.Template{},
		&model.App{},
	)
}
