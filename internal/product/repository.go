package product

import (
	"context"
	"product-service/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, p *domain.Product) (*domain.Product, error)
	FindByID(ctx context.Context, id int64) (*domain.Product, error)
	Update(ctx context.Context, p *domain.Product) error
}
