package model

import (
	"time"

	"github.com/4thena-io/abyss/internal/constant"
)

type Project struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"uniqueIndex"`
	Description string `gorm:"not null"`
	Status      constant.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
