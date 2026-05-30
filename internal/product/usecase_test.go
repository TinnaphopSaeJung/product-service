package product

import (
	"context"
	"errors"
	"product-service/internal/apperrors"
	"product-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeProductRepository struct {
	createCalled    bool
	receivedProduct *domain.Product
	createFunc      func(ctx context.Context, p *domain.Product) (*domain.Product, error)
}

func (f *fakeProductRepository) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	f.createCalled = true
	f.receivedProduct = p

	if f.createFunc != nil {
		return f.createFunc(ctx, p)
	}

	return p, nil
}

func TestCreateProduct_WhenRequestIsValid_ShouldCreateProduct(t *testing.T) {
	description := "Mechanical keyboard"
	salePrice := 1290.0

	repo := &fakeProductRepository{
		createFunc: func(ctx context.Context, p *domain.Product) (*domain.Product, error) {
			return &domain.Product{
				ID:          1,
				Name:        p.Name,
				Description: p.Description,
				SalePrice:   p.SalePrice,
				Price:       p.Price,
			}, nil
		},
	}

	uc := NewUsecase(repo)

	req := CreateProductRequest{
		Name:        "Keyboard",
		Description: &description,
		SalePrice:   &salePrice,
		Price:       1590,
	}

	result, err := uc.CreateProduct(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, repo.createCalled)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, "Keyboard", result.Name)
	assert.Equal(t, &description, result.Description)
	assert.Equal(t, &salePrice, result.SalePrice)
	assert.Equal(t, 1590.0, result.Price)
}

func TestCreateProduct_WhenNameHasLeadingAndTrailingSpaces_ShouldTrimNameBeforeCreate(t *testing.T) {
	repo := &fakeProductRepository{
		createFunc: func(ctx context.Context, p *domain.Product) (*domain.Product, error) {
			return &domain.Product{
				ID:          1,
				Name:        p.Name,
				Description: p.Description,
				SalePrice:   p.SalePrice,
				Price:       p.Price,
			}, nil
		},
	}

	uc := NewUsecase(repo)

	req := CreateProductRequest{
		Name:  " เสื้อ ",
		Price: 100,
	}

	result, err := uc.CreateProduct(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, repo.createCalled)
	assert.Equal(t, "เสื้อ", repo.receivedProduct.Name)
	assert.Equal(t, "เสื้อ", result.Name)
}

func TestCreateProduct_WhenRequestIsInvalid_ShouldReturnValidationErrorAndNotCallRepository(t *testing.T) {
	repo := &fakeProductRepository{}
	uc := NewUsecase(repo)

	req := CreateProductRequest{
		Name:  "",
		Price: 1590,
	}

	result, err := uc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, apperrors.ErrValidation)
	assert.False(t, repo.createCalled)
}

func TestCreateProduct_WhenSalePriceGreaterThanPrice_ShouldReturnValidationErrorAndNotCallRepository(t *testing.T) {
	salePrice := 2000.0

	repo := &fakeProductRepository{}
	uc := NewUsecase(repo)

	req := CreateProductRequest{
		Name:      "Keyboard",
		SalePrice: &salePrice,
		Price:     1590,
	}

	result, err := uc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, apperrors.ErrValidation)
	assert.False(t, repo.createCalled)
}

func TestCreateProduct_WhenRepositoryReturnsError_ShouldReturnError(t *testing.T) {
	repoErr := errors.New("database error")

	repo := &fakeProductRepository{
		createFunc: func(ctx context.Context, p *domain.Product) (*domain.Product, error) {
			return nil, repoErr
		},
	}

	uc := NewUsecase(repo)

	req := CreateProductRequest{
		Name:  "Keyboard",
		Price: 1590,
	}

	result, err := uc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	assert.True(t, repo.createCalled)
}
