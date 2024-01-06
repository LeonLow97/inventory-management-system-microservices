\c imsdb;

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    customer_name VARCHAR(255),
    brand_name VARCHAR(100) NOT NULL,
    category_name VARCHAR(100) NOT NULL,
    color VARCHAR(50),
    size VARCHAR(50),
    quantity INT NOT NULL,
    description TEXT,
    -- https://stackoverflow.com/questions/15726535/which-datatype-should-be-used-for-currency
    -- BIGINT to store as cents, then we can convert to dollars later with highest precision
    revenue BIGINT NOT NULL,
    cost BIGINT NOT NULL,
    profit BIGINT NOT NULL,
    has_reviewed TINYINT(1) DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
