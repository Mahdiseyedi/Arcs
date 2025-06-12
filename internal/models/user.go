package models

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Balance   float64 `gorm:"type:decimal(10,2);not null"`
}
