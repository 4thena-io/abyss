package app

import (
	"time"

	"git.d4ramirez.com/project-abyss/abys-api/internal/constant"
)

type App struct {
	Id          string          `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name" gorm:"not null,uniqueIndex"`
	Description string          `json:"description" gorm:"not null"`
	Kind        string          `json:"kind"`
	Language    string          `json:"language"`
	RepoId      int64           `json:"repo_id" gorm:"not null"`
	RepoUrl     string          `json:"repo_url" gorm:"not null"`
	CloneUrl    string          `json:"clone_url" gorm:"not null"`
	CiId        int64           `json:"ci_id" gorm:"not null"`
	CiUrl       string          `json:"ci_url" gorm:"not null"`
	Project     string          `json:"project" gorm:"not null,foreignKey"`
	Status      constant.Status `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
