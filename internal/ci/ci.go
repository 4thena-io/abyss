package ci

import (
	"context"

	"git.d4ramirez.com/project-abyss/abys-api/internal/dto"
)

type Ci interface {
	ActivateRepo(ctx context.Context, req dto.ActivateRepoRequestDTO) (*dto.ActivateRepoResponseDTO, error)
} 
