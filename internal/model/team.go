package model

import "time"

type Team struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"uniqueIndex"`
	Description string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TeamMember struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	TeamID    uint      `gorm:"not null;index"`
	Team      Team      `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Role      string    `gorm:"not null;default:member"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
