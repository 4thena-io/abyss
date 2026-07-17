package model

import "time"

type Deployment struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	AppID       uint      `gorm:"not null;index"`
	App         App       `gorm:"foreignKey:AppID;constraint:OnDelete:CASCADE"`
	Environment string    `gorm:"not null"`
	Status      string    `gorm:"not null"`
	Commit      string    `gorm:"not null"`
	TriggeredBy string    `gorm:"not null"`
	Duration    int64
	DeployedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
