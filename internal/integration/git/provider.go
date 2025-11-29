package git

import (
	"context"
)

type Provider interface {
	CreateRepo(ctx context.Context, orgName, name string) (*CreateRepoResponseDTO, error)
}

