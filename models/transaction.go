package models

type Transaction struct {
	Op    string  `json:"op"`
	Coin  string  `json:"coin"`
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
	TotalPrice float64 `json:"total_price"`
}
