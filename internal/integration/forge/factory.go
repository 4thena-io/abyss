package forge

import (
	"fmt"

	"git.4thena.io/4thena/abys/internal/config"
	"git.4thena.io/4thena/abys/internal/integration/forge/forgejo"
	"git.4thena.io/4thena/abys/internal/integration/forge/gitea"
	"git.4thena.io/4thena/abys/internal/integration/forge/github"
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
