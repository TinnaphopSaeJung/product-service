package postgres

import (
	"context"
	"errors"
	"product-service/internal/apperrors"
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
	assert.Equal(t, "products_sale_price_lt_price", pgErr.ConstraintName)
}

func TestProductRepository_Create_WhenSalePriceEqualPrice_ShouldReturnConstraintError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	p := &domain.Product{
		Name:      "Invalid Product",
		Price:     1000,
		SalePrice: float64Ptr(1000),
	}

	created, err := repo.Create(context.Background(), p)

	assert.Nil(t, created)
	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	assert.Equal(t, "products_sale_price_lt_price", pgErr.ConstraintName)
}

func TestProductRepository_FindByID_WhenProductExists_ShouldReturnProduct(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	created, err := repo.Create(context.Background(), &domain.Product{
		Name:        "Keyboard",
		Description: stringPtr("Mechanical keyboard"),
		SalePrice:   float64Ptr(1290),
		Price:       1590,
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	found, err := repo.FindByID(context.Background(), created.ID)

	require.NoError(t, err)
	require.NotNil(t, found)

	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "Keyboard", found.Name)
	require.NotNil(t, found.Description)
	assert.Equal(t, "Mechanical keyboard", *found.Description)
	require.NotNil(t, found.SalePrice)
	assert.Equal(t, 1290.0, *found.SalePrice)
	assert.Equal(t, 1590.0, found.Price)
	assert.False(t, found.CreatedAt.IsZero())
	assert.False(t, found.UpdatedAt.IsZero())
}

func TestProductRepository_FindByID_WhenProductDoesNotExist_ShouldReturnProductNotFound(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	found, err := repo.FindByID(context.Background(), 999999)

	assert.Nil(t, found)
	assert.ErrorIs(t, err, apperrors.ErrProductNotFound)
}

func TestProductRepository_Update_WhenProductExists_ShouldUpdateProduct(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	created, err := repo.Create(context.Background(), &domain.Product{
		Name:        "Keyboard",
		Description: stringPtr("Old description"),
		SalePrice:   float64Ptr(1290),
		Price:       1590,
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	created.Name = "Gaming Keyboard"
	created.Description = nil
	created.SalePrice = nil
	created.Price = 1990

	err = repo.Update(context.Background(), created)

	require.NoError(t, err)

	updated, err := repo.FindByID(context.Background(), created.ID)

	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, "Gaming Keyboard", updated.Name)
	assert.Nil(t, updated.Description)
	assert.Nil(t, updated.SalePrice)
	assert.Equal(t, 1990.0, updated.Price)
	assert.False(t, updated.CreatedAt.IsZero())
	assert.False(t, updated.UpdatedAt.IsZero())
}

func TestProductRepository_Update_WhenProductDoesNotExist_ShouldReturnProductNotFound(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	err := repo.Update(context.Background(), &domain.Product{
		ID:          999999,
		Name:        "Not Found Product",
		Description: nil,
		SalePrice:   nil,
		Price:       1000,
	})

	assert.ErrorIs(t, err, apperrors.ErrProductNotFound)
}

func TestProductRepository_Update_WhenSalePriceEqualPrice_ShouldReturnConstraintError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	created, err := repo.Create(context.Background(), &domain.Product{
		Name:      "Keyboard",
		SalePrice: float64Ptr(1290),
		Price:     1590,
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	created.SalePrice = float64Ptr(1590)
	created.Price = 1590

	err = repo.Update(context.Background(), created)

	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	assert.Equal(t, "products_sale_price_lt_price", pgErr.ConstraintName)
}

func TestProductRepository_Update_WhenSalePriceGreaterThanPrice_ShouldReturnConstraintError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	created, err := repo.Create(context.Background(), &domain.Product{
		Name:      "Keyboard",
		SalePrice: float64Ptr(1290),
		Price:     1590,
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	created.SalePrice = float64Ptr(2000)
	created.Price = 1590

	err = repo.Update(context.Background(), created)

	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	assert.Equal(t, "products_sale_price_lt_price", pgErr.ConstraintName)
}

func TestProductRepository_Update_WhenPriceIsZero_ShouldReturnConstraintError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	repo := NewProductRepository(db)

	created, err := repo.Create(context.Background(), &domain.Product{
		Name:  "Keyboard",
		Price: 1590,
	})
	require.NoError(t, err)
	require.NotNil(t, created)

	created.Price = 0

	err = repo.Update(context.Background(), created)

	require.Error(t, err)

	var pgErr *pgconn.PgError
	require.True(t, errors.As(err, &pgErr))
	assert.Equal(t, "products_price_positive", pgErr.ConstraintName)
}
