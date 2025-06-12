package models

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Balance   int64 `gorm:"not null"`
}
