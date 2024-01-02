USE imsdb;

-- brands table
CREATE TABLE brands (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- categories table
CREATE TABLE categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- products table
CREATE TABLE products (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL, -- comes from authentication microservice
    brand_id INT NOT NULL,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    size VARCHAR(50),
    color VARCHAR(50),
    quantity INT NOT NULL,
    is_deleted TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Sample first names
INSERT INTO brands (name) VALUES ('Nike'), ('Adidas'), ('Puma'), ('Balenciaga');

-- Sample last names
INSERT INTO categories (name) VALUES ('Shoes'), ('Clothing'), ('Accessories');

-- Sample product names and descriptions
INSERT INTO products (user_id, brand_id, category_id, name, description, size, color, quantity)
VALUES
    (1, 1, 1, 'Running Shoes', 'Perfect for jogging and running activities', 'Medium', 'Black', 20),
    (1, 2, 2, 'T-Shirt', 'Cotton t-shirt for casual wear', 'Large', 'Blue', 15),
    (2, 3, 3, 'Backpack', 'Spacious backpack for daily use', '', 'Red', 10),
    (3, 1, 1, 'Sneakers', 'Lightweight sneakers for everyday use', 'Small', 'White', 25),
    (2, 3, 2, 'Hoodie', 'Warm hoodie for cold weather', 'Medium', 'Gray', 12),
    (1, 2, 3, 'Sunglasses', 'UV protection sunglasses', '', 'Black', 8),
    (3, 3, 1, 'Walking Shoes', 'Comfortable shoes for walking', 'Large', 'Brown', 18),
    (2, 1, 2, 'Shorts', 'Casual shorts for outdoor activities', 'Small', 'Green', 20),
    (1, 2, 3, 'Watch', 'Stylish wristwatch', '', 'Silver', 5),
    (3, 1, 2, 'Polo Shirt', 'Classic polo shirt for a smart-casual look', 'Medium', 'Red', 15);

