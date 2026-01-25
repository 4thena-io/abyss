package model

import (
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
)

type App struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null,uniqueIndex"`
	Description string `gorm:"not null"`
	Kind        string
	Language    string
	RepoID      int64  `gorm:"not null"`
	RepoURL     string `gorm:"not null"`
	CloneURL    string `gorm:"not null"`
	CIID        int64  `gorm:"not null"`
	CIURL       string `gorm:"not null"`
	ProjectID   uint
	TemplateID  uint `gorm:"not null"`
	Status      constant.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
