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
	findByIDCalled  bool
	updateCalled    bool
	receivedProduct *domain.Product
	receivedID      int64

	createFunc   func(ctx context.Context, p *domain.Product) (*domain.Product, error)
	findByIDFunc func(ctx context.Context, id int64) (*domain.Product, error)
	updateFunc   func(ctx context.Context, p *domain.Product) error
}

func (f *fakeProductRepository) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	f.createCalled = true
	f.receivedProduct = p

	if f.createFunc != nil {
		return f.createFunc(ctx, p)
	}

	return p, nil
}

func (f *fakeProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	f.findByIDCalled = true
	f.receivedID = id

	if f.findByIDFunc != nil {
		return f.findByIDFunc(ctx, id)
	}

	return nil, apperrors.ErrProductNotFound
}

func (f *fakeProductRepository) Update(ctx context.Context, p *domain.Product) error {
	f.updateCalled = true
	f.receivedProduct = p

	if f.updateFunc != nil {
		return f.updateFunc(ctx, p)
	}

	return nil
}

func stringPtr(v string) *string {
	return &v
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

func TestCreateProduct_WhenSalePriceEqualPrice_ShouldReturnValidationErrorAndNotCallRepository(t *testing.T) {
	salePrice := 1590.0

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

func TestCreateProduct_WhenDescriptionIsEmpty_ShouldNormalizeDescriptionToNil(t *testing.T) {
	description := " "

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
		Price:       1590,
	}

	result, err := uc.CreateProduct(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, repo.createCalled)
	assert.Nil(t, repo.receivedProduct.Description)
	assert.Nil(t, result.Description)
}

func TestPatchProduct_WhenPatchName_ShouldUpdateNameOnly(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:          id,
				Name:        "Keyboard",
				Description: stringPtr("Old description"),
				SalePrice:   float64Ptr(1290),
				Price:       1590,
			}, nil
		},
		updateFunc: func(ctx context.Context, p *domain.Product) error {
			return nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Name: OptionalString{
			IsSet: true,
			Value: stringPtr("Gaming Keyboard"),
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	require.NoError(t, err)

	assert.True(t, repo.findByIDCalled)
	assert.True(t, repo.updateCalled)
	assert.Equal(t, int64(1), repo.receivedID)
	assert.Equal(t, "Gaming Keyboard", repo.receivedProduct.Name)
	assert.Equal(t, "Old description", *repo.receivedProduct.Description)
	assert.Equal(t, 1290.0, *repo.receivedProduct.SalePrice)
	assert.Equal(t, 1590.0, repo.receivedProduct.Price)
}

func TestPatchProduct_WhenPatchDescriptionNull_ShouldSetDescriptionToNil(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:          id,
				Name:        "Keyboard",
				Description: stringPtr("Old description"),
				SalePrice:   float64Ptr(1290),
				Price:       1590,
			}, nil
		},
		updateFunc: func(ctx context.Context, p *domain.Product) error {
			return nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Description: OptionalString{
			IsSet: true,
			Value: nil,
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	require.NoError(t, err)

	assert.True(t, repo.findByIDCalled)
	assert.True(t, repo.updateCalled)
	assert.Nil(t, repo.receivedProduct.Description)
	assert.Equal(t, "Keyboard", repo.receivedProduct.Name)
	assert.Equal(t, 1290.0, *repo.receivedProduct.SalePrice)
	assert.Equal(t, 1590.0, repo.receivedProduct.Price)
}

func TestPatchProduct_WhenPatchDescriptionEmpty_ShouldNormalizeDescriptionToNil(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:          id,
				Name:        "Keyboard",
				Description: stringPtr("Old description"),
				SalePrice:   float64Ptr(1290),
				Price:       1590,
			}, nil
		},
		updateFunc: func(ctx context.Context, p *domain.Product) error {
			return nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Description: OptionalString{
			IsSet: true,
			Value: stringPtr("   "),
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	require.NoError(t, err)

	assert.True(t, repo.updateCalled)
	assert.Nil(t, repo.receivedProduct.Description)
}

func TestPatchProduct_WhenPatchSalePriceNull_ShouldSetSalePriceToNil(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:          id,
				Name:        "Keyboard",
				Description: stringPtr("Old description"),
				SalePrice:   float64Ptr(1290),
				Price:       1590,
			}, nil
		},
		updateFunc: func(ctx context.Context, p *domain.Product) error {
			return nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		SalePrice: OptionalFloat64{
			IsSet: true,
			Value: nil,
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	require.NoError(t, err)

	assert.True(t, repo.updateCalled)
	assert.Nil(t, repo.receivedProduct.SalePrice)
	assert.Equal(t, "Keyboard", repo.receivedProduct.Name)
	assert.Equal(t, "Old description", *repo.receivedProduct.Description)
	assert.Equal(t, 1590.0, repo.receivedProduct.Price)
}

func TestPatchProduct_WhenPatchPriceOnlyAndExistingSalePriceBecomesInvalid_ShouldReturnValidationError(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:        id,
				Name:      "Keyboard",
				SalePrice: float64Ptr(1290),
				Price:     1590,
			}, nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Price: OptionalFloat64{
			IsSet: true,
			Value: float64Ptr(1000),
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
	assert.True(t, repo.findByIDCalled)
	assert.False(t, repo.updateCalled)
}

func TestPatchProduct_WhenPatchSalePriceEqualPrice_ShouldReturnValidationError(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:    id,
				Name:  "Keyboard",
				Price: 1590,
			}, nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		SalePrice: OptionalFloat64{
			IsSet: true,
			Value: float64Ptr(1590),
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
	assert.True(t, repo.findByIDCalled)
	assert.False(t, repo.updateCalled)
}

func TestPatchProduct_WhenRequestHasNoField_ShouldReturnValidationErrorAndNotCallRepository(t *testing.T) {
	repo := &fakeProductRepository{}
	uc := NewUsecase(repo)

	req := PatchProductRequest{}

	err := uc.PatchProduct(context.Background(), 1, req)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
	assert.False(t, repo.findByIDCalled)
	assert.False(t, repo.updateCalled)
}

func TestPatchProduct_WhenNameIsNull_ShouldReturnValidationErrorAndNotCallUpdate(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:    id,
				Name:  "Keyboard",
				Price: 1590,
			}, nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Name: OptionalString{
			IsSet: true,
			Value: nil,
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
	assert.True(t, repo.findByIDCalled)
	assert.False(t, repo.updateCalled)
}

func TestPatchProduct_WhenPriceIsNull_ShouldReturnValidationErrorAndNotCallUpdate(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:    id,
				Name:  "Keyboard",
				Price: 1590,
			}, nil
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Price: OptionalFloat64{
			IsSet: true,
			Value: nil,
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
	assert.True(t, repo.findByIDCalled)
	assert.False(t, repo.updateCalled)
}

func TestPatchProduct_WhenProductNotFound_ShouldReturnProductNotFound(t *testing.T) {
	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return nil, apperrors.ErrProductNotFound
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Name: OptionalString{
			IsSet: true,
			Value: stringPtr("New Name"),
		},
	}

	err := uc.PatchProduct(context.Background(), 999, req)

	assert.ErrorIs(t, err, apperrors.ErrProductNotFound)
	assert.True(t, repo.findByIDCalled)
	assert.False(t, repo.updateCalled)
}

func TestPatchProduct_WhenRepositoryUpdateReturnsError_ShouldReturnError(t *testing.T) {
	repoErr := errors.New("database error")

	repo := &fakeProductRepository{
		findByIDFunc: func(ctx context.Context, id int64) (*domain.Product, error) {
			return &domain.Product{
				ID:    id,
				Name:  "Keyboard",
				Price: 1590,
			}, nil
		},
		updateFunc: func(ctx context.Context, p *domain.Product) error {
			return repoErr
		},
	}

	uc := NewUsecase(repo)

	req := PatchProductRequest{
		Name: OptionalString{
			IsSet: true,
			Value: stringPtr("Gaming Keyboard"),
		},
	}

	err := uc.PatchProduct(context.Background(), 1, req)

	assert.ErrorIs(t, err, repoErr)
	assert.True(t, repo.findByIDCalled)
	assert.True(t, repo.updateCalled)
}
