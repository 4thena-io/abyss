package release

import (
	"time"

	"git.d4ramirez.com/project-abyss/abys-api/internal/constant"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (repository *Repository) Save(version *CreateReleaseRecordDto) (*Release, error) {
	newRelease := Release{
		Id:        uuid.NewString(),
		Tag:       version.Tag,
		App:       version.App,
		Status:    constant.StatusNew,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result := repository.db.Create(newRelease)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newRelease, nil
}

func (repository *Repository) GetAllByApp(name string) ([]Release, error) {
	var releases []Release
	result := repository.db.Where("app = ?", name).Find(&releases)
	if result.Error != nil {
		return nil, result.Error
	}
	return releases, nil
}

func (repository *Repository) GetByTag(tag string, name string) (*Release, error) {
	var release Release
	result := repository.db.Where("tag = ?", tag).Where("app = ?", name).First(&release)
	if result.Error != nil {
		return nil, result.Error
	}
	return &release, nil
}
