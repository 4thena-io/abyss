package ci

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/integration/ci/drone"
	"github.com/4thena-io/abyss/internal/integration/ci/gitea_actions"
	"github.com/4thena-io/abyss/internal/integration/ci/woodpecker"
)

func NewCi(cfg config.CIConfig, forgeCfg config.ForgeConfig) (CI, error) {
	host := cfg.Host
	token := cfg.Token

	switch cfg.Type {
	case "gitea-actions", "github-actions", "gitlab-ci":
		if host == "" {
			host = forgeCfg.Host
		}
		if token == "" {
			token = forgeCfg.Token
		}
	}

	switch cfg.Type {
	case "gitea-actions":
		return gitea_actions.NewGiteaActionsCI(host, token)
	case "woodpecker":
		return woodpecker.NewWoodpeckerCi(host, token)
	case "drone":
		return drone.NewDroneCi(host, token)
	default:
		return nil, fmt.Errorf("unknown ci provider %s", cfg.Type)
	}
}
