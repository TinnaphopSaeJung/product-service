package product

import (
	"errors"
	"log"
	"net/http"
	"product-service/internal/apperrors"
	"product-service/internal/response"
	"strconv"

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

func (h *Handler) PatchProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		log.Printf("invalid product id: %v", c.Param("id"))

		c.JSON(http.StatusBadRequest, response.ErrorNoData(response.CodeInvalidProductID))
		return
	}

	var req PatchProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("failed to bind patch product request: %v", err)

		c.JSON(http.StatusBadRequest, response.ErrorNoData(response.CodeValidationError))
		return
	}

	if err = h.usecase.PatchProduct(c.Request.Context(), id, req); err != nil {
		if errors.Is(err, apperrors.ErrValidation) {
			log.Printf("patch product validation error: %v", err)

			c.JSON(http.StatusBadRequest, response.ErrorNoData(response.CodeValidationError))
			return
		}

		if errors.Is(err, apperrors.ErrProductNotFound) {
			log.Printf("patch product not found: id=%d", id)

			c.JSON(http.StatusNotFound, response.ErrorNoData(response.CodeProductNotFound))
			return
		}

		log.Printf("patch product internal error: %v", err)

		c.JSON(http.StatusInternalServerError, response.ErrorNoData(response.CodeInternalServerError))
		return
	}

	c.JSON(http.StatusOK, response.SuccessNoData())
}
