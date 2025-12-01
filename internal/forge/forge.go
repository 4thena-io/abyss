package forge

import (
	"context"

	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
)

type Forge interface {
	CreateRepo(ctx context.Context, owner, name string) (*dto.CreateRepoResponseDTO, error)
	DeleteRepo(ctx context.Context, owner, name string) (error)
}

