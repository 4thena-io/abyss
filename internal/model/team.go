package model

import (
	"time"

	"github.com/4thena-io/abyss/internal/constant"
)

type Team struct {
	ID          uint            `gorm:"primaryKey;autoIncrement"`
	Name        string          `gorm:"uniqueIndex"`
	Description string          `gorm:"not null"`
	Status      constant.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TeamMember struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	TeamID    uint      `gorm:"not null;index"`
	UserID    uint      `gorm:"not null;index"`
	User      User
	Role      string    `gorm:"not null;default:member"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
