package forge

import (
	"fmt"

	"git.4thena.io/4thena/abys/internal/config"
	"git.4thena.io/4thena/abys/internal/integration/forge/gitea"
	"git.4thena.io/4thena/abys/internal/integration/forge/github"
)

var (
	kind	= config.Environment.ForgeType
	host	= config.Environment.ForgeHost
	token	= config.Environment.ForgeToken
)

func NewForge() (Forge, error) {
	switch kind {
		case "github":
			return github.NewGithubForge(host, token)
		case "gitea":
			return gitea.NewGiteaForge(host, token)
		default:
			return nil, fmt.Errorf("unknown forge type %s", kind) 
	}
}
