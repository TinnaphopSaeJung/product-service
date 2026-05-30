package product

import (
	"context"
	"product-service/internal/domain"
	"strings"
)

type Usecase interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*ProductResponse, error)
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
		Description: req.Description,
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
