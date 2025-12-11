package model

import (
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
)

type Release struct {
	Id        string          `json:"id" gorm:"primaryKey"`
	Tag       string          `json:"tag" gorm:"uniqueIndex:uq_version_tag_app"`
	App       string          `json:"app" gorm:"foreignKey;uniqueIndex:uq_version_tag_app"`
	Status    constant.Status `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
