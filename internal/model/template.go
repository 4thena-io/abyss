package model

import (
	"time"

	"github.com/4thena-io/abyss/internal/constant"
)

type Template struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	Kind        string
	Language    string
	RepoURL     string
	CloneURL    string
	Status      constant.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
