package database

import (
	"git.d4ramirez.com/project-abyss/abys-api/internal/model"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.App{},
		&model.Project{},
		&model.Release{},
	)
}
