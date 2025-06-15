package dto

import (
	"arcs/internal/models"
	"time"
)

type CreateUserRequest struct {
	Balance int64 `json:"balance" binding:"required"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type ChargeUserBalance struct {
	UserId string `json:"user_id" binding:"required"`
	Amount int64  `json:"balance" binding:"required"`
}

type GetFilteredUserSMSReq struct {
	UserID string
	Filter models.SMSFilter
}

type GetFilteredUserSMSResp struct {
	SMS   []SMSResponse `json:"list_sms"`
	Count int64         `json:"count"`
}

type SMSResponse struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	OrderID     string    `json:"order_id"`
	Destination string    `json:"destination"`
	Status      string    `json:"status"`
}
