package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	DbType        	string
	DbHost          string
	DbPort          string
	DbUser          string
	DbPassword      string
	DbName          string
	ForgeType				string
	ForgeHost       string
	ForgeToken      string
	CiType					string
	CiHost  				string
	CiToken 				string
}

var Environment = initConfig()

func initConfig() Config {
	return Config{
		DbHost:          getEnv("ABYSS_DB_HOST", "127.0.0.1"),
		DbPort:          getEnv("ABYSS_DB_PORT", "5432"),
		DbType:        	 getEnv("ABYSS_DB_TYPE", "postgres"),
		DbUser:          getEnv("ABYSS_DB_USER", "postgres"),
		DbPassword:      getEnv("ABYSS_DB_PASSWORD", ""),
		DbName:          getEnv("ABYSS_DB_NAME", "postgres"),
		ForgeType: 			 getEnv("FORGE_TYPE", "github"),
		ForgeHost:       getEnv("FORGE_HOST", "http://127.0.0.1:3000"),
		ForgeToken:      getEnv("FORGE_TOKEN", ""),
		CiType: 				 getEnv("CI_TYPE", "jenkins"),
		CiHost:  				 getEnv("CI_HOST", "http://127.0.0.1:3000"),
		CiToken: 				 getEnv("CI_TOKEN", ""),
	}
}

func getEnv(key string, fallback string) string {
	value, ok := os.LookupEnv(key)
	if ok {
		return value
	}
	return fallback
}
