package forge

import (
	"context"

	"git.4thena.io/4thena/abys/internal/dto"
)

type Forge interface {
	CreateRepo(ctx context.Context, owner, name string) (*dto.CreateRepoResponseDTO, error)
	DeleteRepo(ctx context.Context, owner, name string) (error)
}

