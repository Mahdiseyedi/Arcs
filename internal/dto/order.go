package dto

type OrderRequest struct {
	UserID       string   `json:"user_id"`
	Content      string   `json:"content"`
	Destinations []string `json:"destinations"`
}
