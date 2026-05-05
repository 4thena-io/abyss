package repository

import (
	"context"
	"errors"

	"github.com/4thena-io/abyss/internal/model"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) Save(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Create(team).Error
}

func (r *TeamRepository) GetAll(ctx context.Context) ([]model.Team, error) {
	var teams []model.Team
	result := r.db.WithContext(ctx).Find(&teams)
	return teams, result.Error
}

func (r *TeamRepository) GetByID(ctx context.Context, id uint) (*model.Team, error) {
	var team model.Team
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&team)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

func (r *TeamRepository) GetByName(ctx context.Context, name string) (*model.Team, error) {
	var team model.Team
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&team)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

func (r *TeamRepository) Delete(ctx context.Context, team *model.Team) error {
	return r.db.WithContext(ctx).Delete(team).Error
}

func (r *TeamRepository) GetMembers(ctx context.Context, teamID uint) ([]model.TeamMember, error) {
	var members []model.TeamMember
	result := r.db.WithContext(ctx).Where("team_id = ?", teamID).Preload("User").Find(&members)
	return members, result.Error
}

func (r *TeamRepository) GetByUserID(ctx context.Context, userID uint) ([]model.TeamMember, error) {
	var members []model.TeamMember
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Preload("User").Find(&members)
	return members, result.Error
}

func (r *TeamRepository) SaveMember(ctx context.Context, member *model.TeamMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *TeamRepository) CountMembers(ctx context.Context, teamID uint) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.TeamMember{}).Where("team_id = ?", teamID).Count(&count)
	return count, result.Error
}
