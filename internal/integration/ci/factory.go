package ci

import (
	"fmt"

	"git.4thena.io/4thena/abys/internal/config"
	"git.4thena.io/4thena/abys/internal/integration/ci/drone"
	"git.4thena.io/4thena/abys/internal/integration/ci/woodpecker"
)

var (
	kind  = config.Environment.CiType
	host  = config.Environment.CiHost
	token = config.Environment.CiToken
)

func NewCi() (CI, error) {
	switch kind {
	case "woodpecker":
		return woodpecker.NewWoodpeckerCi(host, token)
	case "drone":
		return drone.NewDroneCi(host, token)
	default:
		return nil, fmt.Errorf("unknown ci provider %s", kind)
	}
}
