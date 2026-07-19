package config

type Config struct {
	Server   ServerConfig   `koanf:"server"`
	Database DatabaseConfig `koanf:"database"`
	Forge    ForgeConfig    `koanf:"forge"`
	CI       CIConfig       `koanf:"ci"`
	Auth     AuthConfig     `koanf:"auth"`
	Logging  LoggingConfig  `koanf:"logging"`
}

type ServerConfig struct {
	Host        string `koanf:"host"`
	Port        string `koanf:"port"`
	WebhookHost string `koanf:"webhook_host"`
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
	Type   string `koanf:"type"`
	Host   string `koanf:"host"`
	Token  string `koanf:"token"`
	Owner  string `koanf:"owner"`
	Branch string `koanf:"branch"`
}

type CIConfig struct {
	Type  string `koanf:"type"`
	Host  string `koanf:"host"`
	Token string `koanf:"token"`
}

type AuthConfig struct {
	ClientID     string `koanf:"client_id"`
	ClientSecret string `koanf:"client_secret"`
	CallbackURL  string `koanf:"callback_url"`
}

type LoggingConfig struct {
	Level  string `koanf:"level"`
	Format string `koanf:"format"`
	File   string `koanf:"file"`
}
