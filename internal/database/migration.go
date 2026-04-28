package database

import (
	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Project{},
		&model.Template{},
		&model.App{},
		&model.Deployment{},
	)
}
