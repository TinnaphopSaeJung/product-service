package domain

import "time"

type Product struct {
	ID          int64
	Name        string
	Description *string
	SalePrice   *float64
	Price       float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
