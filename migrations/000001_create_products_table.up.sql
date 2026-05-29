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

    CONSTRAINT products_sale_price_lte_price 
        CHECK (sale_price IS NULL OR sale_price <= price)
);