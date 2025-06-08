package product

import "time"

type Product struct {
	ID           int64     `json:"id"`
	ProductName  string    `json:"product_name"`
	ProductStock int64     `json:"product_stock"`
	ProductPrice int64     `json:"product_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
