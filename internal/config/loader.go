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

func abyssHome() string {
	if h := os.Getenv("ABYSS_HOME"); h != "" {
		return h
	}
	return "."
}

// DefaultPath returns the config file path used when none is specified.
// Override with ABYSS_CONFIG or --config; falls back to $ABYSS_HOME/custom/conf/config.yaml.
func DefaultPath() string {
	if p := os.Getenv("ABYSS_CONFIG"); p != "" {
		return p
	}
	return filepath.Join(abyssHome(), "custom", "conf", "config.yaml")
}

// DefaultDataDir returns the root data directory for persistent storage.
// Override with ABYSS_DATA_DIR; falls back to $ABYSS_HOME/custom/data.
func DefaultDataDir() string {
	if d := os.Getenv("ABYSS_DATA_DIR"); d != "" {
		return d
	}
	return filepath.Join(abyssHome(), "custom", "data")
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
