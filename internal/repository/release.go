package repository

import (
	"context"
	"time"

	"git.4thena.io/4thena/abys/internal/constant"
	"git.4thena.io/4thena/abys/internal/dto"
	"git.4thena.io/4thena/abys/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReleaseRepository struct {
	db *gorm.DB
}

func NewReleaseRepository(db *gorm.DB) *ReleaseRepository {
	return &ReleaseRepository{
		db: db,
	}
}

func (r *ReleaseRepository) SaveRelease(ctx context.Context, version *dto.CreateReleaseRecordDto) (*model.Release, error) {
	newRelease := model.Release{
		Id:        uuid.NewString(),
		Tag:       version.Tag,
		App:       version.App,
		Status:    constant.StatusNew,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := r.db.WithContext(ctx).Create(newRelease)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newRelease, nil
}

func (r *ReleaseRepository) GetAllReleasesByApp(ctx context.Context ,name string) ([]model.Release, error) {
	var releases []model.Release
	result := r.db.WithContext(ctx).Where("app = ?", name).Find(&releases)
	if result.Error != nil {
		return nil, result.Error
	}
	return releases, nil
}

func (r *ReleaseRepository) GetReleaseByTag(ctx context.Context, tag string, name string) (*model.Release, error) {
	var release model.Release
	result := r.db.WithContext(ctx).Where("tag = ?", tag).Where("app = ?", name).First(&release)
	if result.Error != nil {
		return nil, result.Error
	}
	return &release, nil
}
