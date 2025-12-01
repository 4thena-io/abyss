package service

import (
	"context"
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
	"git.d4ramirez.com/project-abyss/abys-api/internal/model"
	"git.d4ramirez.com/project-abyss/abys-api/internal/repository"
	"gorm.io/gorm"
)

type ReleaseService struct {
	repository repository.ReleaseRepository
}

func NewReleaseService(repository repository.ReleaseRepository) *ReleaseService {
	return &ReleaseService{repository}
}

func (s *ReleaseService) GetAllReleases(ctx context.Context, appName string) ([]model.Release, error) {
	data, err := s.repository.GetAllReleasesByApp(ctx, appName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ReleaseService) GetReleaseByTag(ctx context.Context, tag, appName string) (*model.Release, error) {
	data, err := s.repository.GetReleaseByTag(ctx, tag, appName)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("version not found")
	}

	return data, nil
}

func (s *ReleaseService) SaveRelease(ctx context.Context, req dto.CreateReleaseRequestDto, appName string) (*model.Release, error) {
	existing, err := s.repository.GetReleaseByTag(ctx, req.Tag, appName)
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("the version %s already exists for the app %s", req.Tag, appName)
	}

	data, err := s.repository.SaveRelease(ctx, &dto.CreateReleaseRecordDto{
		Tag: req.Tag,
		App: appName,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
