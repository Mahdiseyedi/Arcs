package dto

type OrderRequest struct {
	UserID       string   `json:"user_id" binding:"required"`
	Content      string   `json:"content" binding:"required"`
	Destinations []string `json:"destinations" binding:"required"`
}
