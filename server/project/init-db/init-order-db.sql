\c imsdb;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
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
    has_reviewed BOOLEAN DEFAULT FALSE,
    status VARCHAR(50) NOT NULL,
    status_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW() 
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

INSERT INTO orders (
    product_id, user_id, customer_name,brand_name, category_name, color, size, quantity, description, revenue, cost, profit, has_reviewed, status
) VALUES
    (1, 1, 'John Doe', 'Nike', 'Shoes', 'Black', 'Medium', 2, 'Ordered for daily jogging', 2000, 1500, 500, FALSE, 'COMPLETED'),
    (3, 2, 'Alice Johnson', 'Puma', 'Clothing', 'Red', 'Large', 1, 'Bought for casual wear', 800, 600, 200, FALSE, 'FAILED'),
    (6, 3, 'Bob Smith', 'Adidas', 'Accessories', 'Black', '', 1, 'Ordered for UV protection', 150, 100, 50, FALSE, 'SUBMITTED'),
    (4, 1, 'Emily Brown', 'Nike', 'Shoes', 'White', 'Small', 3, 'Purchased for everyday use', 2400, 1800, 600, FALSE, 'COMPLETED'),
    (8, 2, 'David Lee', 'Adidas', 'Clothing', 'Green', 'Small', 2, 'Ordered for outdoor activities', 400, 300, 100, FALSE, 'COMPLETED'),
    (9, 3, 'Sophia Wilson', 'Nike', 'Accessories', 'Silver', '', 1, 'Purchased for a stylish look', 120, 80, 40, FALSE, 'FAILED'),
    (5, 1, 'Olivia Martinez', 'Puma', 'Clothing', 'Gray', 'Medium', 1, 'Ordered for cold weather', 1000, 750, 250, FALSE, 'FAILED'),
    (2, 2, 'William Johnson', 'Adidas', 'Shoes', 'Blue', 'Large', 2, 'Ordered for casual jogging', 1600, 1200, 400, FALSE, 'FAILED'),
    (10, 3, 'Charlotte Davis', 'Nike', 'Clothing', 'Red', 'Medium', 1, 'Purchased for a smart-casual look', 300, 200, 100, FALSE, 'SUBMITTED'),
    (7, 1, 'James Taylor', 'Puma', 'Accessories', 'Brown', '', 1, 'Ordered for daily usage', 200, 150, 50, FALSE, 'SUBMITTED');
