package models

import (
	"time"
)

type Order struct {
	ID           string `gorm:"type:uuid;primary_key"`
	CreatedAt    time.Time
	UserID       string   `gorm:"type:uuid;not null"`
	Content      string   `gorm:"type:text"`
	Destinations []string `gorm:"serializer:json"`
}
