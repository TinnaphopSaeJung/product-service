package product

import (
	"errors"
	"log"
	"net/http"
	"product-service/internal/apperrors"
	"product-service/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase Usecase
}

func NewHandler(usecase Usecase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to bind create product request: %v", err)

		c.JSON(http.StatusBadRequest, response.Error(response.CodeValidationError))
		return
	}

	result, err := h.usecase.CreateProduct(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, apperrors.ErrValidation) {
			log.Printf("create product validation error: %v", err)

			c.JSON(http.StatusBadRequest, response.Error(response.CodeValidationError))
			return
		}

		log.Printf("create product internal server error: %v", err)

		c.JSON(http.StatusInternalServerError, response.Error(response.CodeInternalServerError))
		return
	}

	c.JSON(http.StatusCreated, response.Success(result))
}
