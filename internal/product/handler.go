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

// CreateProduct godoc
// @Summary Create product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param request body CreateProductRequest true "Create product request"
// @Success 201 {object} CreateProductSuccessSwaggerResponse
// @Failure 400 {object} CreateProductErrorSwaggerResponse
// @Failure 500 {object} CreateProductErrorSwaggerResponse
// @Router /product [post]
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

// PatchProduct godoc
// @Summary Patch product
// @Description Partially update product by ID. Only provided fields will be updated.
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param request body PatchProductSwaggerRequest true "Patch product request"
// @Success 200 {object} PatchProductSuccessSwaggerResponse
// @Failure 400 {object} PatchProductErrorSwaggerResponse
// @Failure 404 {object} PatchProductErrorSwaggerResponse
// @Failure 500 {object} PatchProductErrorSwaggerResponse
// @Router /product/{id} [patch]
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
