package model

import "time"

type App struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"not null,uniqueIndex"`
	Description  string `gorm:"not null"`
	Kind         string
	Language     string
	RepoID       int64  `gorm:"not null"`
	RepoFullName string `gorm:"not null"`
	RepoURL      string `gorm:"not null"`
	CloneURL     string `gorm:"not null"`
	CIID         int64  `gorm:"not null"`
	CISlug       string
	CIURL        string `gorm:"not null"`
	ProjectID    uint
	Project      *Project  `gorm:"foreignKey:ProjectID;constraint:OnDelete:SET NULL"`
	TemplateID   *uint
	Template     *Template `gorm:"foreignKey:TemplateID;constraint:OnDelete:SET NULL"`
	CreatorID    uint
	Creator      User      `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
