package model

type Order struct {
	ID        int64  `json:"id"`
	UserEmail string `json:"user_email"`
	ProductID int64  `json:"product_id"`
	Quantity  int    `json:"quantity"`
	Address   string `json:"address"`
	Status    string `json:"status"`
}
