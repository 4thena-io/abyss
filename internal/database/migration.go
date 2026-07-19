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
	{
		// Replaces the plaintext User.Token column with TokenHash/TokenLastEight.
		// Runs AutoMigrate to add the new columns, then drops the old one so no
		// raw PAT is left sitting on disk.
		ID: "20260719000001_hash_personal_access_tokens",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&model.User{}); err != nil {
				return err
			}
			if tx.Migrator().HasColumn(&model.User{}, "token") {
				return tx.Migrator().DropColumn(&model.User{}, "token")
			}
			return nil
		},
		Rollback: func(tx *gorm.DB) error {
			// The raw tokens are gone for good; rollback only removes the hash
			// columns so users regenerate a PAT afterward.
			if err := tx.Migrator().DropColumn(&model.User{}, "token_hash"); err != nil {
				return err
			}
			return tx.Migrator().DropColumn(&model.User{}, "token_last_eight")
		},
	},
}

func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations)
	return m.Migrate()
}
