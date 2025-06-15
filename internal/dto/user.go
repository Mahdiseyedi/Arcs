package dto

import "arcs/internal/models"

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
	SMS   []models.SMS `json:"list_sms"`
	Count int64        `json:"count"`
}
