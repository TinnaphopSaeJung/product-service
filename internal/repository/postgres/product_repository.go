package postgres

import (
	"context"
	"fmt"
	"product-service/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	query := `
		INSERT INTO products (name, description, sale_price, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, sale_price, price, created_at, updated_at
	`

	var created domain.Product

	err := r.db.QueryRow(
		ctx,
		query,
		p.Name,
		p.Description,
		p.SalePrice,
		p.Price,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.SalePrice,
		&created.Price,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &created, nil
}
