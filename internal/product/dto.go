package product

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type CreateProductRequest struct {
	Name        string   `json:"name" example:"Keyboard"`
	Description *string  `json:"description" example:"Mechanical keyboard"`
	SalePrice   *float64 `json:"sale_price" example:"1290"`
	Price       float64  `json:"price" example:"1590"`
}
type ProductResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}

type PatchProductRequest struct {
	Name        OptionalString  `json:"name"`
	Description OptionalString  `json:"description"`
	SalePrice   OptionalFloat64 `json:"sale_price"`
	Price       OptionalFloat64 `json:"price"`
}

func (r PatchProductRequest) HasAnyField() bool {
	return r.Name.IsSet ||
		r.Description.IsSet ||
		r.SalePrice.IsSet ||
		r.Price.IsSet
}

type OptionalString struct {
	IsSet bool
	Value *string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.IsSet = true

	if bytes.Equal(data, []byte("null")) {
		o.Value = nil
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("invalid string value: %w", err)
	}

	o.Value = &value
	return nil
}

type OptionalFloat64 struct {
	IsSet bool
	Value *float64
}

func (o *OptionalFloat64) UnmarshalJSON(data []byte) error {
	o.IsSet = true

	if bytes.Equal(data, []byte("null")) {
		o.Value = nil
		return nil
	}

	var value float64
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("invalid float64 value: %w", err)
	}

	o.Value = &value
	return nil
}
