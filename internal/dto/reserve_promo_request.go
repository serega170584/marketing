package dto

type ReservePromoRequestDTO struct {
	Code    string `json:"code"`
	UserID  string `json:"user_id"`
	OrderID string `json:"order_id"`
}
