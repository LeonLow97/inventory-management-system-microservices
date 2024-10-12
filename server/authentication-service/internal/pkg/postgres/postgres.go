package postgres_conn

import (
	"fmt"
	"log"
	"regexp"

	"github.com/LeonLow97/internal/pkg/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

func ConnectToPostgreSQL(cfg config.Config) (*sqlx.DB, error) {
	// Construct the DSN string based on environment variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DB,
	)

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		errMsg := err.Error()
		errMsg = regexp.MustCompile(`(://[^:]+:).+(@.+)`).ReplaceAllString(errMsg, "$1*****$2")
		errMsg = regexp.MustCompile(`(password=).+(\s+)`).ReplaceAllString(errMsg, "$1*****$2")
		return nil, fmt.Errorf("parsing DSN failed: %s", errMsg)
	}
	connStr := stdlib.RegisterConnConfig(connConfig)
	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		log.Println("failed to open postgres connection with error:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("failed to ping postgres database with error:", err)
		return nil, err
	}

	if cfg.Mode == config.ModeDocker {
		if err := CreateTableAndInsertDummyData(db); err != nil {
			log.Println("failed to create tables and insert dummy data with error:", err)
		}
	}

	log.Println("Successfully connected to Postgres database!")
	return db, nil
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
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			username VARCHAR(50) NOT NULL UNIQUE,
			password VARCHAR(60) NOT NULL,
			email VARCHAR(100) NOT NULL,
			active INT NOT NULL DEFAULT 1,
			admin INT NOT NULL DEFAULT 0,
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		);

		CREATE INDEX idx_users_username ON users (username);

		INSERT INTO users (first_name, last_name, username, password, email)
		VALUES
			('Jie Wei', 'Low', 'lowjiewei', '$2a$10$OULOXURo57bo5keyNXGQxefqMyEM67JIscqLVKWgd/S.siCqNAHdC', 'lowjiewei@email.com'),
			('Leon', 'Low', 'leonlow', '$2a$10$OULOXURo57bo5keyNXGQxefqMyEM67JIscqLVKWgd/S.siCqNAHdC', 'leonlow@email.com');
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
