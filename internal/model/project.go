package model

import "time"

type Project struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"uniqueIndex"`
	Description string `gorm:"not null"`
	TeamID      *uint
	Team        *Team `gorm:"foreignKey:TeamID;constraint:OnDelete:SET NULL"`
	CreatorID   uint
	Creator     User  `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
