package ci

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/integration/ci/drone"
	"github.com/4thena-io/abyss/internal/integration/ci/forgejo_actions"
	"github.com/4thena-io/abyss/internal/integration/ci/gitea_actions"
	"github.com/4thena-io/abyss/internal/integration/ci/github_actions"
	"github.com/4thena-io/abyss/internal/integration/ci/gitlab_ci"
	"github.com/4thena-io/abyss/internal/integration/ci/woodpecker"
)

func NewCi(cfg config.CIConfig, forgeCfg config.ForgeConfig) (CI, error) {
	host := cfg.Host
	token := cfg.Token

	switch cfg.Type {
	case "gitea-actions", "forgejo-actions", "github-actions", "gitlab-ci":
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
	case "forgejo-actions":
		return forgejo_actions.NewForgejoActionsCI(host, token)
	case "github-actions":
		return github_actions.NewGithubActionsCI(host, token)
	case "gitlab-ci":
		return gitlab_ci.NewGitlabCI(host, token)
	case "woodpecker":
		return woodpecker.NewWoodpeckerCi(host, token)
	case "drone":
		return drone.NewDroneCi(host, token)
	default:
		return nil, fmt.Errorf("unknown ci provider %s", cfg.Type)
	}
}
