package product

import (
	"context"
	"fmt"
	"product-service/internal/apperrors"
	"product-service/internal/domain"
	"strings"
)

type Usecase interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
	PatchProduct(ctx context.Context, id int64, req PatchProductRequest) error
}

type usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error) {
	p := &domain.Product{
		Name:        strings.TrimSpace(req.Name),
		Description: normalizeNullableString(req.Description),
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	if err := ValidateProduct(*p); err != nil {
		return nil, err
	}

	createdProduct, err := u.repo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return &ProductResponse{
		ID:          createdProduct.ID,
		Name:        createdProduct.Name,
		Description: createdProduct.Description,
		SalePrice:   createdProduct.SalePrice,
		Price:       createdProduct.Price,
	}, nil
}

func normalizeNullableString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func (u *usecase) PatchProduct(ctx context.Context, id int64, req PatchProductRequest) error {
	if !req.HasAnyField() {
		return fmt.Errorf("%w: at least one field is required", apperrors.ErrValidation)
	}

	existingProduct, err := u.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if req.Name.IsSet {
		if req.Name.Value == nil {
			return fmt.Errorf("%w: name cannot be null", apperrors.ErrValidation)
		}

		existingProduct.Name = strings.TrimSpace(*req.Name.Value)
	}

	if req.Description.IsSet {
		existingProduct.Description = normalizeNullableString(req.Description.Value)
	}

	if req.SalePrice.IsSet {
		existingProduct.SalePrice = req.SalePrice.Value
	}

	if req.Price.IsSet {
		if req.Price.Value == nil {
			return fmt.Errorf("%w: price cannot be null", apperrors.ErrValidation)
		}

		existingProduct.Price = *req.Price.Value
	}

	if err := ValidateProduct(*existingProduct); err != nil {
		return err
	}

	if err := u.repo.Update(ctx, existingProduct); err != nil {
		return err
	}

	return nil
}
