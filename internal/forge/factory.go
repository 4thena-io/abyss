package forge

import (
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/config"
)

var (
	kind	= config.Environment.ForgeType
	host	= config.Environment.ForgeHost
	token	= config.Environment.ForgeToken
)

func NewForge() (Forge, error) {
	switch kind {
		case "gitea":
			return NewGiteaForge(host, token)
		default:
			return nil, fmt.Errorf("unknown forge type %s", kind) 
	}
}
