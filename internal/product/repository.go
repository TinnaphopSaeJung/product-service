package product

import (
	"context"
	"product-service/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, p *domain.Product) (*domain.Product, error)
}
