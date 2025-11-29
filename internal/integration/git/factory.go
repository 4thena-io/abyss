package git

import (
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/config"
)

func NewGitProvider(providerType string) (Provider, error) {
	switch providerType {
		case "gitea":
			return NewGiteaProvider(config.Environment.GiteaHost, config.Environment.GiteaToken)
		default:
		return nil, fmt.Errorf("unkwown git ptovider %s", providerType) 
	}
}
