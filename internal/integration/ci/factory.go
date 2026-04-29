package ci

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/integration/ci/drone"
	"github.com/4thena-io/abyss/internal/integration/ci/gitea_actions"
	"github.com/4thena-io/abyss/internal/integration/ci/woodpecker"
)

func effectiveCIHost(cfg config.Config) string {
	if cfg.CiHost != "" {
		return cfg.CiHost
	}

	switch cfg.CiType {
	case "gitea-actions", "github-actions", "gitlab-ci":
		return cfg.ForgeHost
	default:
		return ""
	}
}

func effectiveCIToken(cfg config.Config) string {
	if cfg.CiToken != "" {
		return cfg.CiToken
	}

	switch cfg.CiType {
	case "gitea-actions", "github-actions", "gitlab-ci":
		return cfg.ForgeToken
	default:
		return ""
	}
}

func NewCi() (CI, error) {
	cfg := config.Environment

	host := effectiveCIHost(cfg)
	token := effectiveCIToken(cfg)

	switch cfg.CiType {
	case "gitea-actions":
		return gitea_actions.NewGiteaActionsCI(host, token)

	case "woodpecker":
		return woodpecker.NewWoodpeckerCi(host, token)

	case "drone":
		return drone.NewDroneCi(host, token)

	default:
		return nil, fmt.Errorf("unknown ci provider %s", cfg.CiType)
	}
}
