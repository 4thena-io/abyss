package repository

import (
	"context"
	"errors"

	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

type ServerConfigRepository struct {
	db *gorm.DB
}

func NewServerConfigRepository(db *gorm.DB) *ServerConfigRepository {
	return &ServerConfigRepository{db: db}
}

// Get returns the value for the given key, or ("", nil) if the key does not exist.
func (r *ServerConfigRepository) Get(ctx context.Context, key string) (string, error) {
	var cfg model.ServerConfig
	result := r.db.WithContext(ctx).Where("key = ?", key).First(&cfg)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return "", nil
	}
	return cfg.Value, result.Error
}

// Set upserts the given key/value pair.
func (r *ServerConfigRepository) Set(ctx context.Context, key, value string) error {
	cfg := model.ServerConfig{Key: key, Value: value}
	return r.db.WithContext(ctx).Save(&cfg).Error
}
