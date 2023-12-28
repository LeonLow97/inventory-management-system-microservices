## Inventory Service

### Endpoints

| Method | Endpoint                         | Description                              |
| :----: | -------------------------------- | ---------------------------------------- |
|  GET   | `/inventory/products`            | Retrieve a list of products              |
|  GET   | `/inventory/products/{id}`       | Retrieve details of a specific product   |
|  POST  | `/inventory/products`            | Create a new product                     |
|  PUT   | `/inventory/products/{id}`       | Update details of a product              |
| DELETE | `/inventory/products/{id}`       | Delete a product                         |
|  POST  | `/inventory/products/{id}/adjust | Adjust the inventory count for a product |

### Database Tables

- **Products Table**
    - Stores the product information including name, brand, size, color, quantity, category, and potentially additional details.
- **Brands Table**
    - Contains brand details associated with products. Each product can have a relationship with a brand through a foreign key reference.
- **Categories Table**
    - Holds category information for products. This table establishes the category relationship for each product.