package database

import (
	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.ServerConfig{},
		&model.User{},
		&model.Project{},
		&model.Template{},
		&model.App{},
		&model.Team{},
		&model.TeamMember{},
		&model.Deployment{},
	)
}
