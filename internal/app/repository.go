package app

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

func (repository *Repository) Save(app *CreateAppRecordDto) (*App, error) {
	newApp := App{
		Id:          uuid.NewString(),
		Name:        app.Name,
		Description: app.Description,
		RepoId:      app.RepoId,
		RepoUrl:     app.RepoUrl,
		CloneUrl:    app.CloneUrl,
		CiId:        app.CiId,
		CiUrl:       app.CiUrl,
		Project:     app.Project,
		Status:      constant.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := repository.db.Create(newApp)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newApp, nil
}

func (repository *Repository) GetAll() ([]App, error) {
	var applications []App
	result := repository.db.Find(&applications)
	if result.Error != nil {
		return nil, result.Error
	}
	return applications, nil
}

func (repository *Repository) GetById(id string) (*App, error) {
	var application App
	result := repository.db.Where("id = ?", id).First(&application)
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
}

func (repository *Repository) GetByName(name string) (*App, error) {
	var application App
	result := repository.db.Where("name = ?", name).First(&application)
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
}
