package models

import "time"

type SMS struct {
	CreatedAt   time.Time
	ID          string `gorm:"type:uuid;not null"`
	UserID      string `gorm:"type:uuid;not null"`
	OrderID     string `gorm:"type:uuid;not null"`
	Destination string
	Status      string
}
