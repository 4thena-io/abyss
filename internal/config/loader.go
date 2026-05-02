package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func Load(configPath string) (*Config, error) {
	k := koanf.New(".")

	if err := k.Load(confmap.Provider(defaults(), "."), nil); err != nil {
		return nil, fmt.Errorf("failed to load defaults: %w", err)
	}

	if configPath == "" {
		configPath = DefaultPath()
	}

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
		}
	}

	if err := k.Load(env.Provider("ABYSS__", ".", func(s string) string {
		s = strings.TrimPrefix(s, "ABYSS__")
		s = strings.ToLower(s)
		return strings.ReplaceAll(s, "__", ".")
	}), nil); err != nil {
		return nil, fmt.Errorf("failed to load env vars: %w", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// DefaultPath returns the config file path used when none is specified.
// Set ABYSS_HOME to override the base directory (defaults to the working directory).
func DefaultPath() string {
	base := os.Getenv("ABYSS_HOME")
	if base == "" {
		base = "."
	}
	return filepath.Join(base, "custom", "conf", "config.yaml")
}

func defaults() map[string]interface{} {
	return map[string]interface{}{
		"server.host":    "0.0.0.0",
		"server.port":    "8000",
		"database.type":  "sqlite",
		"database.path":  "abyss.db",
		"database.host":  "127.0.0.1",
		"database.port":  "5432",
		"database.user":  "postgres",
		"database.name":  "abyss",
		"forge.type":      "gitea",
		"forge.host":      "http://127.0.0.1:3000",
		"ci.type":         "gitea-actions",
		"logging.level":   "info",
		"logging.format":  "pretty",
	}
}
