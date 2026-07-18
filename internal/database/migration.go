package database

import (
	"github.com/4thena-io/abyss/internal/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// migrations is the ordered list of schema changes. Append new migrations to
// the bottom of the list with a new ID; never edit or reorder existing ones
// once they've shipped, since gormigrate tracks applied IDs in a
// "migrations" table and only runs the ones it hasn't seen.
var migrations = []*gormigrate.Migration{
	{
		ID: "20260718000001_initial_schema",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(
				&model.ServerConfig{},
				&model.User{},
				&model.Project{},
				&model.Template{},
				&model.App{},
				&model.Team{},
				&model.TeamMember{},
				&model.Deployment{},
			)
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(
				&model.Deployment{},
				&model.TeamMember{},
				&model.Team{},
				&model.App{},
				&model.Template{},
				&model.Project{},
				&model.User{},
				&model.ServerConfig{},
			)
		},
	},
}

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations)
	return m.Migrate()
}
