package postgres

import (
	"context"
	"errors"
	"fmt"
	"product-service/internal/apperrors"
	"product-service/internal/domain"

	"github.com/jackc/pgx/v5"
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

func (r *ProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := `
		SELECT id, name, description, sale_price, price, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var p domain.Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.SalePrice,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrProductNotFound
		}

		return nil, fmt.Errorf("failed to find product by id: %w", err)
	}

	return &p, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product) error {
	query := `
		UPDATE products
		SET
			name = $1,
			description = $2,
			sale_price = $3,
			price = $4,
			updated_at = NOW()
		WHERE id = $5
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		p.Name,
		p.Description,
		p.SalePrice,
		p.Price,
		p.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return apperrors.ErrProductNotFound
	}

	return nil
}
