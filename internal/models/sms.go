package models

import (
	"arcs/internal/models/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type SMS struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	OrderID string `gorm:"type:uuid;not null"`

	Destination string
	Status      string

	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (s *SMS) ToProto() *proto.SMS {
	return &proto.SMS{
		Id:          s.ID,
		CreatedAt:   timestamppb.New(s.CreatedAt),
		UpdatedAt:   timestamppb.New(s.UpdatedAt),
		OrderId:     s.OrderID,
		Destination: s.Destination,
		Status:      s.Status,
	}
}
