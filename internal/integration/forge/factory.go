package forge

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/integration/forge/forgejo"
	"github.com/4thena-io/abyss/internal/integration/forge/gitea"
	"github.com/4thena-io/abyss/internal/integration/forge/github"
	"github.com/4thena-io/abyss/internal/integration/forge/gitlab"
)

func NewForge(cfg config.ForgeConfig) (Forge, error) {
	switch cfg.Type {
	case "gitea":
		return gitea.NewGiteaForge(cfg.Host, cfg.Token)
	case "forgejo":
		return forgejo.NewForgejoForge(cfg.Host, cfg.Token)
	case "github":
		return github.NewGithubForge(cfg.Host, cfg.Token)
	case "gitlab":
		return gitlab.NewGitlabForge(cfg.Host, cfg.Token)
	default:
		return nil, fmt.Errorf("unknown forge type %s", cfg.Type)
	}
}
