package component

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"product-service/internal/app"
	"product-service/internal/testutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type apiResponse struct {
	Successful bool            `json:"successful"`
	ErrorCode  string          `json:"error_code"`
	Data       json.RawMessage `json:"data"`
}

type productResponse struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
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
	assert.Equal(t, "", res.ErrorCode)
	require.NotNil(t, res.Data)

	var product productResponse
	err = json.Unmarshal(res.Data, &product)
	require.NoError(t, err)

	assert.Equal(t, int64(1), product.ID)
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
	assert.Equal(t, "", res.ErrorCode)

	var product productResponse
	err = json.Unmarshal(res.Data, &product)
	require.NoError(t, err)

	assert.Equal(t, int64(1), product.ID)
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
	assert.Equal(t, "VALIDATION_ERROR", res.ErrorCode)
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
	assert.Equal(t, "VALIDATION_ERROR", res.ErrorCode)
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
	assert.Equal(t, "VALIDATION_ERROR", res.ErrorCode)
	assert.Equal(t, "null", string(res.Data))
}
