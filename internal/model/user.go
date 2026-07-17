package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ForgeID   int64   `gorm:"uniqueIndex"`
	Username  string  `gorm:"uniqueIndex"`
	Email     string
	IsAdmin   bool
	AvatarURL string
	Token     *string `gorm:"uniqueIndex"` // personal access token for CLI; NULL when not yet generated
}
