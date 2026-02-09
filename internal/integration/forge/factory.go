package forge

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/integration/forge/forgejo"
	"github.com/4thena-io/abyss/internal/integration/forge/gitea"
	"github.com/4thena-io/abyss/internal/integration/forge/github"
)

var (
	kind  = config.Environment.ForgeType
	host  = config.Environment.ForgeHost
	token = config.Environment.ForgeToken
)

func NewForge() (Forge, error) {
	switch kind {
	case "gitea":
		return gitea.NewGiteaForge(host, token)
	case "forgejo":
		return forgejo.NewForgejoForge(host, token)
	case "github":
		return github.NewGithubForge(host, token)
	default:
		return nil, fmt.Errorf("unknown forge type %s", kind)
	}
}
