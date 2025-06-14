package models

import (
	"time"
)

type Order struct {
	ID        string `gorm:"type:uuid;primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID  string `gorm:"type:uuid;not null"`
	Content string `gorm:"type:text"`

	User *User `gorm:"foreignKey:UserID;references:ID"`
	SMS  []SMS `gorm:"foreignKey:OrderID"`
}
