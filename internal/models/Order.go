package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key"`
	CreatedAt    time.Time
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	Content      string    `gorm:"type:text"`
	Destinations []string  `gorm:"serializer:json"`
}
