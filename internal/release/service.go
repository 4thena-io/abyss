package release

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository}
}

func (s *Service) GetAllReleases(ctx context.Context, appName string) ([]Release, error) {
	data, err := s.repository.GetAllByApp(appName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *Service) GetReleaseByTag(ctx context.Context, tag, appName string) (*Release, error) {
	data, err := s.repository.GetByTag(tag, appName)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("version not found")
	}

	return data, nil
}

func (s *Service) SaveRelease(ctx context.Context, req CreateReleaseRequestDto, appName string) (*Release, error) {
	existing, err := s.repository.GetByTag(req.Tag, appName)
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("the version %s already exists for the app %s", req.Tag, appName)
	}

	//TODO Create git tag

	data, err := s.repository.Save(&CreateReleaseRecordDto{
		Tag: req.Tag,
		App: appName,
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
