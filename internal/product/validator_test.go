package product

import (
	"product-service/internal/apperrors"
	"product-service/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func float64Ptr(v float64) *float64 {
	return &v
}

func TestValidateProduct_WhenProductIsValid_ShouldReturnNil(t *testing.T) {
	p := domain.Product{
		Name:      "Keyboard",
		Price:     1590,
		SalePrice: float64Ptr(1290),
	}

	err := ValidateProduct(p)

	assert.NoError(t, err)
}

func TestValidateProduct_WhenNameIsEmpty_ShouldReturnValidationError(t *testing.T) {
	p := domain.Product{
		Name:  "",
		Price: 1590,
	}

	err := ValidateProduct(p)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestValidateProduct_WhenNameHasOnlySpaces_ShouldReturnValidationError(t *testing.T) {
	p := domain.Product{
		Name:  "   ",
		Price: 1590,
	}

	err := ValidateProduct(p)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestValidateProduct_WhenPriceIsZero_ShouldReturnValidationError(t *testing.T) {
	p := domain.Product{
		Name:  "Keyboard",
		Price: 0,
	}

	err := ValidateProduct(p)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestValidateProduct_WhenPriceIsNegative_ShouldReturnValidationError(t *testing.T) {
	p := domain.Product{
		Name:  "Keyboard",
		Price: -100,
	}

	err := ValidateProduct(p)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestValidateProduct_WhenSalePriceIsNegative_ShouldReturnValidationError(t *testing.T) {
	p := domain.Product{
		Name:      "Keyboard",
		Price:     1590,
		SalePrice: float64Ptr(-100),
	}

	err := ValidateProduct(p)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestValidateProduct_WhenSalePriceGreaterThanPrice_ShouldReturnValidationError(t *testing.T) {
	p := domain.Product{
		Name:      "Keyboard",
		Price:     1590,
		SalePrice: float64Ptr(2000),
	}

	err := ValidateProduct(p)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
}

func TestValidateProduct_WhenSalePriceIsNil_ShouldReturnNil(t *testing.T) {
	p := domain.Product{
		Name:      "Keyboard",
		Price:     1590,
		SalePrice: nil,
	}

	err := ValidateProduct(p)

	assert.NoError(t, err)
}

func TestValidateProduct_WhenSalePriceEqualPrice_ShouldReturnValidationError(t *testing.T) {
	p := domain.Product{
		Name:      "Keyboard",
		Price:     1590,
		SalePrice: float64Ptr(1590),
	}

	err := ValidateProduct(p)

	assert.ErrorIs(t, err, apperrors.ErrValidation)
}
