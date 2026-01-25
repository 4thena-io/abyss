package model

import "time"

type Config struct {
	ID                  uint `gorm:"primaryKey;autoIncrement"`
	ServiceAccountName  string
	ServiceAccountEmail string
	ServiceAccountToken string
	ForgeURL            string
	CIURL               string
	CIType              string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
