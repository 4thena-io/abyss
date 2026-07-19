package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ForgeID   int64  `gorm:"uniqueIndex"`
	Username  string `gorm:"uniqueIndex"`
	Email     string
	IsAdmin   bool
	AvatarURL string
	// TokenHash/TokenLastEight replace storing the raw PAT: the raw value is
	// shown once at generation time and never persisted. Lookup goes through
	// the indexed last-eight chars to find the candidate row, then a
	// constant-time compare against the salted hash confirms the match.
	TokenHash      *string `gorm:"index"`
	TokenLastEight string  `gorm:"index"`
}
