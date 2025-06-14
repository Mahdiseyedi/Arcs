package models

import (
	"time"
)

type User struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Balance   int64 `gorm:"not null"`

	Orders []Order `gorm:"foreignKey:UserID"`
}
