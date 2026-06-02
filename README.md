# Product Service

Product Service is a Golang REST API service for managing product information.

This project was implemented using Clean Architecture principles, dependency injection, PostgreSQL, Swagger API documentation, Docker, and automated tests across multiple levels.

## Tech Stack

* Golang
* Gin Web Framework
* PostgreSQL
* pgx
* Swagger / Swaggo
* Docker
* Docker Compose
* Testify

## Features

* Create product
* Patch product by ID
* Partial update support
* Nullable field handling
* Business validation
* PostgreSQL database constraints
* Swagger API documentation
* Unit tests
* Usecase unit tests
* Repository integration tests
* Component tests / E2E within service

## Project Structure

```text
product-service/
├── cmd/
│   └── api/
│       └── main.go                         # Application entrypoint
│
├── docker/
│   └── postgres/
│       └── init.sql                        # PostgreSQL initialization script for Docker Compose
│
├── docs/
│   ├── docs.go                             # Generated Swagger docs
│   ├── swagger.json                        # Generated Swagger JSON file
│   └── swagger.yaml                        # Generated Swagger YAML file
│
├── internal/
│   ├── app/
│   │   └── router.go                       # Gin router setup and dependency injection
│   │
│   ├── apperrors/
│   │   └── errors.go                       # Application-level reusable errors
│   │
│   ├── config/
│   │   └── config.go                       # Environment configuration loader
│   │
│   ├── database/
│   │   └── postgres.go                     # PostgreSQL connection pool setup
│   │
│   ├── domain/
│   │   └── product.go                      # Product domain model
│   │
│   ├── product/
│   │   ├── dto.go                          # Request/response DTOs and PATCH optional types
│   │   ├── handler.go                      # HTTP handlers
│   │   ├── repository.go                   # Repository interface
│   │   ├── routes.go                       # Product route registration
│   │   ├── swagger_dto.go                      # Swagger-only request/response models
│   │   ├── usecase.go                      # Product business usecases
│   │   ├── usecase_test.go                 # Usecase unit tests
│   │   ├── validator.go                    # Product business validation
│   │   └── validator_test.go               # Service/domain validation tests
│   │
│   ├── repository/
│   │   └── postgres/
│   │       ├── product_repository.go        # PostgreSQL repository implementation
│   │       └── product_repository_test.go   # Repository integration tests
│   │
│   ├── response/
│   │   └── response.go                     # Standard API response format
│   │
│   └── testutil/
│       └── db.go                           # Test database helper
│
├── migrations/
│   ├── 000001_create_products_table.up.sql     # SQL migration for creating products table
│   └── 000001_create_products_table.down.sql   # SQL migration rollback
│
├── tests/
│   └── component/
│       └── product_api_test.go             # Component tests / E2E within service
│
├── .dockerignore
├── .env.example
├── .gitignore
├── Dockerfile                              # API Docker image build
├── docker-compose.yml                      # Local Docker environment
├── go.mod
├── go.sum
└── README.md
```

## API Documentation

Swagger UI is available at:

```text
http://localhost:8080/api-docs/index.html
```

The Swagger route is registered under:

```text
/api-docs/*
```

Generated Swagger files are located in:

```text
docs/
```

To regenerate Swagger documentation:

```bash
swag init -g cmd/api/main.go
```

## API Endpoints

### Health Check

```http
GET /health
```

#### Success Response

```json
{
  "successful": true,
  "error_code": null,
  "data": {
    "status": "ok"
  }
}
```

---

## Create Product

```http
POST /product
```

### Request Body

```json
{
  "name": "Keyboard",
  "description": "Mechanical keyboard",
  "sale_price": 1290,
  "price": 1590
}
```

### Field Rules

| Field       | Type   | Required | Nullable | Rule                                                   |
| ----------- | ------ | -------: | -------: | ------------------------------------------------------ |
| name        | string |      Yes |       No | Must not be empty                                      |
| description | string |       No |      Yes | Empty string will be stored as NULL                    |
| sale_price  | number |       No |      Yes | Must be greater than or equal to 0 and less than price |
| price       | number |      Yes |       No | Must be greater than 0                                 |

### Success Response

```http
201 Created
```

```json
{
  "successful": true,
  "error_code": null,
  "data": {
    "id": 1,
    "name": "Keyboard",
    "description": "Mechanical keyboard",
    "sale_price": 1290,
    "price": 1590
  }
}
```

### Error Response

```http
400 Bad Request
```

```json
{
  "successful": false,
  "error_code": "VALIDATION_ERROR",
  "data": null
}
```

---

## Patch Product

```http
PATCH /product/{id}
```

This endpoint supports partial update.

Only fields sent in the request body will be updated.

### Request Body Example

```json
{
  "name": "Gaming Keyboard"
}
```

### Nullable Update Examples

The following request clears `description` and stores it as `NULL`.

```json
{
  "description": null
}
```

The following request clears `sale_price` and stores it as `NULL`.

```json
{
  "sale_price": null
}
```

### Field Behavior

| Field       | Undefined           | Null     | Value              |
| ----------- | ------------------- | -------- | ------------------ |
| name        | Keep existing value | Invalid  | Update name        |
| description | Keep existing value | Set NULL | Update description |
| sale_price  | Keep existing value | Set NULL | Update sale_price  |
| price       | Keep existing value | Invalid  | Update price       |

### Business Rules

* `name` must not be empty.
* `price` must be greater than 0.
* `sale_price` must be greater than or equal to 0.
* `sale_price` must be less than `price`.
* At least one field must be provided in PATCH request body.

### Success Response

```http
200 OK
```

```json
{
  "successful": true,
  "error_code": null
}
```

### Error Responses

#### Invalid Product ID

```http
400 Bad Request
```

```json
{
  "successful": false,
  "error_code": "INVALID_PRODUCT_ID"
}
```

#### Validation Error

```http
400 Bad Request
```

```json
{
  "successful": false,
  "error_code": "VALIDATION_ERROR"
}
```

#### Product Not Found

```http
404 Not Found
```

```json
{
  "successful": false,
  "error_code": "PRODUCT_NOT_FOUND"
}
```

---

## Error Codes

| Error Code            | Description                                    |
| --------------------- | ---------------------------------------------- |
| VALIDATION_ERROR      | Request body or business validation is invalid |
| INVALID_PRODUCT_ID    | Product ID path parameter is invalid           |
| PRODUCT_NOT_FOUND     | Product does not exist                         |
| INTERNAL_SERVER_ERROR | Unexpected server error                        |

## Run with Docker Compose

### Start Application

```bash
docker compose up --build
```

This command starts:

* PostgreSQL
* Product Service API

API will be available at:

```text
http://localhost:8080
```

Swagger UI will be available at:

```text
http://localhost:8080/api-docs/index.html
```

PostgreSQL will be available from host machine at:

```text
localhost:5433
```

PostgreSQL inside Docker network is available at:

```text
postgres:5432
```

### Stop Application

```bash
docker compose down
```

### Stop Application and Remove Database Volume

```bash
docker compose down -v
```

> Note: `docker compose down -v` removes the PostgreSQL data volume.

## Run Tests with Docker

Run all tests inside a Docker container:

```bash
docker compose --profile test run --rm test
```

This command starts the `test` service defined in `docker-compose.yml`.

The `test` service runs the following command inside the container:

```bash
go test ./... -v -count=1 -p 1
```

Explanation:

* `./...` runs tests from all packages.
* `-v` shows detailed test output.
* `-count=1` disables Go test cache and forces tests to run again.
* `-p 1` runs test packages sequentially to avoid database test conflicts.

The test service connects to the `product_test` database using `TEST_DB_*` environment variables.

## Run Locally

### 1. Create `.env`

Copy `.env.example` to `.env`.

```bash
cp .env.example .env
```

Example configuration when using PostgreSQL from Docker Compose:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=1234
DB_NAME=product_db
DB_SSLMODE=disable

TEST_DB_HOST=localhost
TEST_DB_PORT=5433
TEST_DB_USER=postgres
TEST_DB_PASSWORD=1234
TEST_DB_NAME=product_test
TEST_DB_SSLMODE=disable
```

### 2. Start PostgreSQL

```bash
docker compose up postgres
```

### 3. Run Application Locally

```bash
go run cmd/api/main.go
```

### 4. Run Tests Locally

```bash
go test ./... -v -count=1
```

## Database

The service uses PostgreSQL.

Main database:

```text
product_db
```

Test database:

```text
product_test
```

The Docker PostgreSQL initialization script is located at:

```text
docker/postgres/init.sql
```

Migration files are located at:

```text
migrations/
```

### Products Table

```sql
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NULL,
    sale_price NUMERIC(12,2) NULL,
    price NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT products_price_positive 
        CHECK (price > 0),

    CONSTRAINT products_sale_price_non_negative 
        CHECK (sale_price IS NULL OR sale_price >= 0),

    CONSTRAINT products_sale_price_lt_price 
        CHECK (sale_price IS NULL OR sale_price < price)
);
```

## Test Strategy

This project includes 4 levels of testing.

### 1. Service / Domain Unit Test

File:

```text
internal/product/validator_test.go
```

Purpose:

* Tests product business validation.
* Example rules:

  * `name` is required.
  * `price` must be greater than 0.
  * `sale_price` must be greater than or equal to 0.
  * `sale_price` must be less than `price`.

### 2. Usecase Unit Test

File:

```text
internal/product/usecase_test.go
```

Purpose:

* Tests usecase orchestration.
* Uses fake repository.
* Does not connect to a real database.
* Tests create and patch product business flow.

### 3. Repository Integration Test

File:

```text
internal/repository/postgres/product_repository_test.go
```

Purpose:

* Tests real repository implementation with PostgreSQL.
* Verifies:

  * insert product
  * find product by ID
  * update product
  * database constraints

### 4. Component Test / E2E Within Service

File:

```text
tests/component/product_api_test.go
```

Purpose:

* Tests HTTP request through Gin router.
* Covers full flow inside the service:

```text
HTTP Request
→ Handler
→ Usecase
→ Repository
→ PostgreSQL Test DB
→ HTTP Response
```

## Useful Commands

### Run All Tests Locally

```bash
go test ./... -v -count=1
```

### Run All Tests with Docker

```bash
docker compose --profile test run --rm test
```

### Run Only Product Unit Tests

```bash
go test ./internal/product -v -count=1
```

### Run Only Repository Integration Tests

```bash
go test ./internal/repository/postgres -v -count=1
```

### Run Only Component Tests

```bash
go test ./tests/component -v -count=1
```

### Check Running Containers

```bash
docker compose ps
```

### View API Logs

```bash
docker compose logs -f api
```

### Connect to PostgreSQL Container

```bash
docker exec -it product_service_postgres psql -U postgres -d product_db
```

### Query Products

```sql
SELECT * FROM products ORDER BY id ASC;
```

### Show NULL Values Clearly in psql

```sql
\pset null '[NULL]'
```

## Notes

* `POST /product` returns product data in the `data` field.
* `PATCH /product/{id}` returns no `data` field according to the API specification.
* `description` and `sale_price` are nullable fields.
* Empty `description` is normalized to `NULL`.
* `sale_price` must always be less than `price`.
* Application-level validation is implemented in Go.
* Database-level validation is enforced using PostgreSQL CHECK constraints.
* Swagger documentation is available at `/api-docs/index.html`.
