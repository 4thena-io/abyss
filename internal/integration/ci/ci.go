package ci

import (
	"context"

	"git.4thena.io/4thena/abys/internal/dto/request"
	"git.4thena.io/4thena/abys/internal/dto/response"
)

type Ci interface {
	ActivateRepo(ctx context.Context, req request.ActivateRepo) (*response.ActivateRepo, error)
	GetBuilds(ctx context.Context, repoID int64) ([]response.Build, error)
} 
