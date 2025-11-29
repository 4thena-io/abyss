package project

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

func (repository *Repository) Save(project *CreateProjectRecordDto) (*Project, error) {
	newProject := Project{
		Id:          uuid.NewString(),
		Name:        project.Name,
		Description: project.Description,
		Status:      constant.StatusNew,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := repository.db.Create(newProject)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newProject, nil
}

func (repository *Repository) GetAll() ([]Project, error) {
	var projects []Project
	result := repository.db.Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}

func (repository *Repository) GetById(id string) (*Project, error) {
	var project Project
	result := repository.db.Where("id = ?", id).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}

func (repository *Repository) GetByName(name string) (*Project, error) {
	var project Project
	result := repository.db.Where("name = ?", name).First(&project)
	if result.Error != nil {
		return nil, result.Error
	}
	return &project, nil
}
