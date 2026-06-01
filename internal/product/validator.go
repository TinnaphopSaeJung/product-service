package product

import (
	"fmt"
	"product-service/internal/apperrors"
	"product-service/internal/domain"
	"strings"
)

func ValidateProduct(p domain.Product) error {
	if strings.TrimSpace(p.Name) == "" {
		return fmt.Errorf("%w: name is required", apperrors.ErrValidation)
	}

	if p.Price <= 0 {
		return fmt.Errorf("%w: price must be greater than zero", apperrors.ErrValidation)
	}

	if p.SalePrice != nil {
		if *p.SalePrice < 0 {
			return fmt.Errorf("%w: sale_price must be greater than or equal to zero", apperrors.ErrValidation)
		}

		if *p.SalePrice >= p.Price {
			return fmt.Errorf("%w: sale_price must be less than price", apperrors.ErrValidation)
		}
	}

	return nil
}
