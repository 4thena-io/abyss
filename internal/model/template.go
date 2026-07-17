package model

import "time"

type Template struct {
	ID          uint `gorm:"primaryKey;autoIncrement"`
	Name        string
	Description string
	Kind        string
	Language    string
	RepoURL     string
	CloneURL    string
	CreatorID   uint
	Creator     User `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
