package model

type ServerConfig struct {
	Key   string `gorm:"primaryKey"`
	Value string  `gorm:"not null"`
}
