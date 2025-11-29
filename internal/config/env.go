package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Address         string
	DbSchema        string
	DbHost          string
	DbPort          string
	DbUser          string
	DbPassword      string
	DbName          string
	GiteaHost       string
	GiteaToken      string
	WoodpeckerHost  string
	WoodpeckerToken string
}

var Environment = initConfig()

func initConfig() Config {
	return Config{
		Address:         getEnv("ABYSS_ADDR", "http://127.0.0.1:8000"),
		DbHost:          getEnv("ABYSS_DB_HOST", "127.0.0.1"),
		DbPort:          getEnv("ABYSS_DB_PORT", "5432"),
		DbSchema:        getEnv("ABYSS_DB_SCHEMA", "public"),
		DbUser:          getEnv("ABYSS_DB_USER", "postgres"),
		DbPassword:      getEnv("ABYSS_DB_PASSWORD", ""),
		DbName:          getEnv("ABYSS_DB_NAME", "postgres"),
		GiteaHost:       getEnv("GITEA_HOST", "http://127.0.0.1:3000"),
		GiteaToken:      getEnv("GITEA_TOKEN", ""),
		WoodpeckerHost:  getEnv("WOODPECKER_HOST", "http://127.0.0.1:3000"),
		WoodpeckerToken: getEnv("WOODPECKER_TOKEN", ""),
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallback
}
