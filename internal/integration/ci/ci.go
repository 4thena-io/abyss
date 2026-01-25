package ci

import (
	"context"

	"git.4thena.io/4thena/abys/internal/model"
)

type CI interface {
	ActivateRepo(ctx context.Context, ID int64) (*model.CIRepo, error)
	GetBuilds(ctx context.Context, ID int64) ([]model.Build, error)
}
