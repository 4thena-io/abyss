package ci

import (
	"fmt"

	"git.d4ramirez.com/project-abyss/abys-api/internal/config"
)

var (
	kind = config.Environment.CiType
	host = config.Environment.CiHost
	token = config.Environment.CiToken
)

func NewCi() (Ci, error) {
	switch kind {
		case "woodpecker":
			return NewWoodpeckerCi(host, token)
		default:
			return nil, fmt.Errorf("unkwown ci ptovider %s", kind) 
	}
}
