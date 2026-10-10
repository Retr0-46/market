CREATE TABLE product_types (
    id BIGSERIAL PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,      -- 'sport', 'electronics', ...
    name TEXT NOT NULL              -- отображаемое имя, например 'Спортивные товары'
);

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    type_id BIGINT NOT NULL REFERENCES product_types(id),
    seller_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    price BIGINT NOT NULL CHECK (price >= 0),
    stock_quantity BIGINT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0)
);

CREATE INDEX idx_products_type_id ON products(type_id);