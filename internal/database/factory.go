package database

import (
	"fmt"
)

func NewDatabaseProvider(providerType string) (Provider, error) {
	switch providerType {
	case "postgres":
		return NewPostgresProvider(), nil
	case "mysql":
		return NewMySQLProvider(), nil
	default:
		return nil, fmt.Errorf("unknown database provider: %s", providerType)
	}
}
