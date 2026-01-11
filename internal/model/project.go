package model

import (
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
)

type Project struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"uniqueIndex"`
	Description string `gorm:"not null"`
	Status      constant.Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
