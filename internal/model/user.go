package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ForgeID   int64  `gorm:"uniqueIndex"`
	Username  string `gorm:"uniqueIndex"`
	IsAdmin   bool
	AvatarURL string
	Token     *string `gorm:"uniqueIndex"` // personal access token for CLI; NULL when not yet generated
}
