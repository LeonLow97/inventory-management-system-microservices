package postgres_conn

import (
	"fmt"
	"log"

	"github.com/LeonLow97/internal/pkg/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectToDB(cfg config.Config) *sqlx.DB {
	// Construct the DSN string based on environment variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DB,
	)
	log.Printf("Connecting to database with URL: %s", dsn)

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		log.Fatal("failed to parse config with error:", err)
		return nil
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Fatal("failed to open connection with pgx with error:", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping postgres database with error:", err)
		return nil
	}

	if cfg.Mode == config.ModeDocker {
		if err := CreateTableAndInsertDummyData(db); err != nil {
			log.Println("failed to create tables and insert dummy data with error:", err)
		}
	}

	log.Println("Successfully connected to Postgres database!")
	return db
}

// CreateTableAndInsertDummyData creates a table and inserts dummy data into it.
func CreateTableAndInsertDummyData(db *sqlx.DB) error {
	// Use a transaction to combine the operations
	tx, err := db.Beginx()
	if err != nil {
		log.Println("failed to begin transaction:", err)
		return err
	}

	// Define the combined SQL statement
	sql := `
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
	`

	// Execute the combined SQL statement
	if _, err := tx.Exec(sql); err != nil {
		log.Println("failed to execute SQL statements:", err)
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Println("failed to rollback transaction:", rollbackErr)
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Println("failed to commit transaction:", err)
		return err
	}

	log.Println("Successfully created table and inserted dummy data!")
	return nil
}
