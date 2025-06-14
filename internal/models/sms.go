package models

import "time"

type SMS struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	OrderID string `gorm:"type:uuid;not null"`

	Destination string
	Status      string

	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}
