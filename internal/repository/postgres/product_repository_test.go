package postgres

import (
	"context"
	"errors"
	"product-service/internal/domain"
	"product-service/internal/testutil"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringPtr(v string) *string {
	return &v
}

func float64Ptr(v float64) *float64 {
	return &v
}

func TestProductRepository_Create_WhenValidProduct_ShouldInsertProduct(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	p := &domain.Product{
		Name:        "Keyboard",
		Description: stringPtr("Mechanical keyboard"),
		SalePrice:   float64Ptr(1290),
		Price:       1590,
	}

	created, err := repo.Create(context.Background(), p)

	require.NoError(t, err)
	require.NotNil(t, created)

	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, "Keyboard", created.Name)
	assert.Equal(t, "Mechanical keyboard", *created.Description)
	assert.Equal(t, 1290.0, *created.SalePrice)
	assert.Equal(t, 1590.0, created.Price)
	assert.False(t, created.CreatedAt.IsZero())
	assert.False(t, created.UpdatedAt.IsZero())
}

func TestProductRepository_Create_WhenNullableFieldsAreNil_ShouldInsertNull(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	p := &domain.Product{
		Name:        "Mouse",
		Description: nil,
		SalePrice:   nil,
		Price:       790,
	}

	created, err := repo.Create(context.Background(), p)

	require.NoError(t, err)
	require.NotNil(t, created)

	assert.Equal(t, int64(1), created.ID)
	assert.Equal(t, "Mouse", created.Name)
	assert.Nil(t, created.Description)
	assert.Nil(t, created.SalePrice)
	assert.Equal(t, 790.0, created.Price)
}

func TestProductRepository_Create_WhenPriceIsZero_ShouldReturnConstraintError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	p := &domain.Product{
		Name:  "Invalid Product",
		Price: 0,
	}

	created, err := repo.Create(context.Background(), p)

	assert.Nil(t, created)
	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	assert.Equal(t, "products_price_positive", pgErr.ConstraintName)
}

func TestProductRepository_Create_WhenSalePriceGreaterThanPrice_ShouldReturnConstraintError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	p := &domain.Product{
		Name:      "Invalid Product",
		Price:     1000,
		SalePrice: float64Ptr(1500),
	}

	created, err := repo.Create(context.Background(), p)

	assert.Nil(t, created)
	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	assert.Equal(t, "products_sale_price_lte_price", pgErr.ConstraintName)
}
