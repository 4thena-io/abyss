package ci

import (
	"context"

	"git.4thena.io/4thena/abys/internal/dto"
)

type Ci interface {
	ActivateRepo(ctx context.Context, req dto.ActivateRepoRequestDTO) (*dto.ActivateRepoResponseDTO, error)
} 
