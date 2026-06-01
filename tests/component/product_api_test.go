package component

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"product-service/internal/app"
	postgresRepo "product-service/internal/repository/postgres"
	"product-service/internal/testutil"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type apiResponse struct {
	Successful bool            `json:"successful"`
	ErrorCode  *string         `json:"error_code"`
	Data       json.RawMessage `json:"data"`
}

type productResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}

func createProductForTest(t *testing.T, router *gin.Engine, body string) productResponse {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	require.True(t, res.Successful)
	require.Nil(t, res.ErrorCode)
	require.NotNil(t, res.Data)

	var product productResponse
	err = json.Unmarshal(res.Data, &product)
	require.NoError(t, err)

	return product
}

func TestPostProduct_WhenRequestIsValid_ShouldReturnCreatedProduct(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "Keyboard",
		"description": "Mechanical keyboard",
		"sale_price": 1290,
		"price": 1590
	}`)

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.True(t, res.Successful)
	assert.Nil(t, res.ErrorCode)
	require.NotNil(t, res.Data)

	var product productResponse
	err = json.Unmarshal(res.Data, &product)
	require.NoError(t, err)

	assert.Greater(t, product.ID, int64(0))
	assert.Equal(t, "Keyboard", product.Name)
	require.NotNil(t, product.Description)
	assert.Equal(t, "Mechanical keyboard", *product.Description)
	require.NotNil(t, product.SalePrice)
	assert.Equal(t, 1290.0, *product.SalePrice)
	assert.Equal(t, 1590.0, product.Price)
}

func TestPostProduct_WhenNullableFieldsAreNull_ShouldReturnCreatedProductWithNullFields(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "Mouse",
		"description": null,
		"sale_price": null,
		"price": 790
	}`)

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.True(t, res.Successful)
	assert.Nil(t, res.ErrorCode)
	require.NotNil(t, res.Data)

	var product productResponse
	err = json.Unmarshal(res.Data, &product)
	require.NoError(t, err)

	assert.Greater(t, product.ID, int64(0))
	assert.Equal(t, "Mouse", product.Name)
	assert.Nil(t, product.Description)
	assert.Nil(t, product.SalePrice)
	assert.Equal(t, 790.0, product.Price)
}

func TestPostProduct_WhenNameIsEmpty_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "",
		"price": 1590
	}`)

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Equal(t, "null", string(res.Data))
}

func TestPostProduct_WhenSalePriceGreaterThanPrice_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "Keyboard",
		"sale_price": 2000,
		"price": 1590
	}`)

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Equal(t, "null", string(res.Data))
}

func TestPostProduct_WhenSalePriceEqualPrice_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "Keyboard",
		"sale_price": 1590,
		"price": 1590
	}`)

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Equal(t, "null", string(res.Data))
}

func TestPostProduct_WhenJsonBodyIsInvalid_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "Keyboard",
		"price":
	}`)

	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Equal(t, "null", string(res.Data))
}

func TestPatchProduct_WhenPatchName_ShouldReturnSuccessNoDataAndUpdateProduct(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"description": "Old description",
		"sale_price": 1290,
		"price": 1590
	}`)

	body := []byte(`{
		"name": "Gaming Keyboard"
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.True(t, res.Successful)
	assert.Nil(t, res.ErrorCode)
	assert.Nil(t, res.Data)

	repo := postgresRepo.NewProductRepository(db)

	updated, err := repo.FindByID(t.Context(), created.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, "Gaming Keyboard", updated.Name)
	require.NotNil(t, updated.Description)
	assert.Equal(t, "Old description", *updated.Description)
	require.NotNil(t, updated.SalePrice)
	assert.Equal(t, 1290.0, *updated.SalePrice)
	assert.Equal(t, 1590.0, updated.Price)
}

func TestPatchProduct_WhenPatchDescriptionNull_ShouldReturnSuccessNoDataAndClearDescription(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"description": "Old description",
		"sale_price": 1290,
		"price": 1590
	}`)

	body := []byte(`{
		"description": null
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.True(t, res.Successful)
	assert.Nil(t, res.ErrorCode)
	assert.Nil(t, res.Data)

	repo := postgresRepo.NewProductRepository(db)

	updated, err := repo.FindByID(t.Context(), created.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, "Keyboard", updated.Name)
	assert.Nil(t, updated.Description)
	require.NotNil(t, updated.SalePrice)
	assert.Equal(t, 1290.0, *updated.SalePrice)
	assert.Equal(t, 1590.0, updated.Price)
}

func TestPatchProduct_WhenPatchDescriptionEmpty_ShouldReturnSuccessNoDataAndClearDescription(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"description": "Old description",
		"sale_price": 1290,
		"price": 1590
	}`)

	body := []byte(`{
		"description": "   "
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.True(t, res.Successful)
	assert.Nil(t, res.ErrorCode)
	assert.Nil(t, res.Data)

	repo := postgresRepo.NewProductRepository(db)

	updated, err := repo.FindByID(t.Context(), created.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Nil(t, updated.Description)
}

func TestPatchProduct_WhenPatchSalePriceNull_ShouldReturnSuccessNoDataAndClearSalePrice(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"description": "Old description",
		"sale_price": 1290,
		"price": 1590
	}`)

	body := []byte(`{
		"sale_price": null
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.True(t, res.Successful)
	assert.Nil(t, res.ErrorCode)
	assert.Nil(t, res.Data)

	repo := postgresRepo.NewProductRepository(db)

	updated, err := repo.FindByID(t.Context(), created.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, "Keyboard", updated.Name)
	require.NotNil(t, updated.Description)
	assert.Equal(t, "Old description", *updated.Description)
	assert.Nil(t, updated.SalePrice)
	assert.Equal(t, 1590.0, updated.Price)
}

func TestPatchProduct_WhenPatchPriceOnlyMakesExistingSalePriceInvalid_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"sale_price": 1290,
		"price": 1590
	}`)

	body := []byte(`{
		"price": 1000
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Nil(t, res.Data)

	repo := postgresRepo.NewProductRepository(db)

	found, err := repo.FindByID(t.Context(), created.ID)
	require.NoError(t, err)
	require.NotNil(t, found)

	assert.Equal(t, 1590.0, found.Price)
	require.NotNil(t, found.SalePrice)
	assert.Equal(t, 1290.0, *found.SalePrice)
}

func TestPatchProduct_WhenSalePriceEqualPrice_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"price": 1590
	}`)

	body := []byte(`{
		"sale_price": 1590
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Nil(t, res.Data)
}

func TestPatchProduct_WhenRequestHasNoField_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"price": 1590
	}`)

	body := []byte(`{}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Nil(t, res.Data)
}

func TestPatchProduct_WhenNameIsNull_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"price": 1590
	}`)

	body := []byte(`{
		"name": null
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Nil(t, res.Data)
}

func TestPatchProduct_WhenPriceIsNull_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	created := createProductForTest(t, router, `{
		"name": "Keyboard",
		"price": 1590
	}`)

	body := []byte(`{
		"price": null
	}`)

	req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/product/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Nil(t, res.Data)
}

func TestPatchProduct_WhenProductDoesNotExist_ShouldReturnProductNotFound(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "New Name"
	}`)

	req := httptest.NewRequest(http.MethodPatch, "/product/999999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "PRODUCT_NOT_FOUND", *res.ErrorCode)
	assert.Nil(t, res.Data)
}

func TestPatchProduct_WhenProductIDIsInvalid_ShouldReturnInvalidProductID(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name": "New Name"
	}`)

	req := httptest.NewRequest(http.MethodPatch, "/product/abc", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "INVALID_PRODUCT_ID", *res.ErrorCode)
	assert.Nil(t, res.Data)
}

func TestPatchProduct_WhenJsonBodyIsInvalid_ShouldReturnValidationError(t *testing.T) {
	db := testutil.NewTestPostgresPool(t)
	defer db.Close()

	testutil.CleanupProducts(t, db)

	router := app.NewRouter(db)

	body := []byte(`{
		"name":
	}`)

	req := httptest.NewRequest(http.MethodPatch, "/product/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)

	var res apiResponse
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	require.NoError(t, err)

	assert.False(t, res.Successful)
	require.NotNil(t, res.ErrorCode)
	assert.Equal(t, "VALIDATION_ERROR", *res.ErrorCode)
	assert.Nil(t, res.Data)
}
