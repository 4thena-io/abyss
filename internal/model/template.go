package model

import (
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
)

type Template struct {
	Id          string          `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Kind        string          `json:"kind"`
	Language    string          `json:"language"`
	RepoUrl     string          `json:"repo_url"`
	CloneUrl    string          `json:"clone_url"`
	Status      constant.Status `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
