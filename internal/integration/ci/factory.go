package ci

import (
	"fmt"

	"github.com/4thena-io/abyss/internal/config"
	"github.com/4thena-io/abyss/internal/integration/ci/drone"
	"github.com/4thena-io/abyss/internal/integration/ci/woodpecker"
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
