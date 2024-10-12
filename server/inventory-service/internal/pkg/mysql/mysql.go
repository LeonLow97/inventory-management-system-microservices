package mysql_conn

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/LeonLow97/internal/pkg/config"
)

func ConnectToMySQL(cfg config.Config) *sql.DB {
	// MySQL DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.MySQLConfig.User,
		cfg.MySQLConfig.Password,
		cfg.MySQLConfig.Host,
		cfg.MySQLConfig.Port,
		cfg.MySQLConfig.Database,
	)

	// open a connection to MySQL database
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening connection to mysql database in inventory service", err)
	}

	// ping mysql database
	if err = conn.Ping(); err != nil {
		log.Fatal("Error pinging mysql database", err)
	}

	if cfg.Mode == config.ModeDocker {
		if err := CreateTableAndInsertDummyData(conn); err != nil {
			log.Println("failed to create tables and insert dummy data with error:", err)
		}
	}

	log.Println("Successfully connected to MySQL database!")
	return conn
}

// CreateTableAndInsertDummyData creates tables and inserts dummy data into them for MySQL.
func CreateTableAndInsertDummyData(db *sql.DB) error {
	// Use a transaction to combine the operations
	tx, err := db.Begin()
	if err != nil {
		log.Println("failed to begin transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Println("failed to rollback transaction:", rollbackErr)
			}
		}
	}()

	// Define the SQL statements for table creation and data insertion
	statements := []string{
		`CREATE TABLE IF NOT EXISTS brands (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS products (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			brand_id INT NOT NULL,
			category_id INT NOT NULL,
			name VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			size VARCHAR(50),
			color VARCHAR(50),
			quantity INT NOT NULL,
			is_deleted TINYINT(1) DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (brand_id) REFERENCES brands(id) ON DELETE CASCADE,
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
		);`,
		`INSERT INTO brands (name) VALUES ('Nike'), ('Adidas'), ('Puma'), ('Balenciaga');`,
		`INSERT INTO categories (name) VALUES ('Shoes'), ('Clothing'), ('Accessories');`,
		`INSERT INTO products (user_id, brand_id, category_id, name, description, size, color, quantity)
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
			(3, 1, 2, 'Polo Shirt', 'Classic polo shirt for a smart-casual look', 'Medium', 'Red', 15);`,
	}

	// Execute each SQL statement in the transaction
	for _, stmt := range statements {
		_, err := tx.Exec(stmt)
		if err != nil {
			log.Println("failed to execute SQL statement:", err)
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Println("failed to commit transaction:", err)
		return err
	}

	log.Println("Successfully created tables and inserted dummy data!")
	return nil
}
