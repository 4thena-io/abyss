package model

import (
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
)

type Project struct {
	Id          string          `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name" gorm:"uniqueIndex"`
	Description string          `json:"description" gorm:"not null"`
	Status      constant.Status `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
