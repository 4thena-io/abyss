package repository

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
	"gorm.io/gorm"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{
		db: db,
	}
}

func (r *ConfigRepository) Save(ctx context.Context, config model.Config) error {
	result := r.db.WithContext(ctx).Create(config)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
