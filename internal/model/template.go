package model

import (
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
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
