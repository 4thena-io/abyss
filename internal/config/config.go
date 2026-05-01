package config

type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	Forge    ForgeConfig    `koanf:"forge"`
	CI       CIConfig       `koanf:"ci"`
	Auth     AuthConfig     `koanf:"auth"`
}

type ServerConfig struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}

type DatabaseConfig struct {
	Type     string `koanf:"type"`
	Path     string `koanf:"path"`
	Host     string `koanf:"host"`
	Port     string `koanf:"port"`
	User     string `koanf:"user"`
	Password string `koanf:"password"`
	Name     string `koanf:"name"`
}

type ForgeConfig struct {
	Type  string `koanf:"type"`
	Host  string `koanf:"host"`
	Token string `koanf:"token"`
	Owner string `koanf:"owner"`
}

type CIConfig struct {
	Type  string `koanf:"type"`
	Host  string `koanf:"host"`
	Token string `koanf:"token"`
}

type AuthConfig struct {
	ClientID      string `koanf:"client_id"`
	ClientSecret  string `koanf:"client_secret"`
	SessionSecret string `koanf:"session_secret"`
}
