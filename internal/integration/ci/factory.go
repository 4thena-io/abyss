package ci

import (
	"fmt"

	"git.4thena.io/4thena/abys/internal/config"
	"git.4thena.io/4thena/abys/internal/integration/ci/jenkins"
	"git.4thena.io/4thena/abys/internal/integration/ci/woodpecker"
)

var (
	kind = config.Environment.CiType
	host = config.Environment.CiHost
	token = config.Environment.CiToken
)

func NewCi() (CI, error) {
	switch kind {
		case "jenkins":
			return jenkins.NewJenkinsCi(host, token)
		case "woodpecker":
			return woodpecker.NewWoodpeckerCi(host, token)
		default:
			return nil, fmt.Errorf("unkwown ci ptovider %s", kind) 
	}
}
