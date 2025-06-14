package dto

type CreateUserRequest struct {
	Balance int64 `json:"balance"`
}

type ChargeUserBalance struct {
	UserId string `json:"user_id"`
	Amount int64  `json:"balance"`
}
