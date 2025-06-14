package dto

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
