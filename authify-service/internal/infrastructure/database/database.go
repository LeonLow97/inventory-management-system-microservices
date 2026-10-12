package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// DB wraps sqlx.DB.
type DB struct {
	*sqlx.DB
}

var (
	instance *DB
	once     sync.Once
)

// GetDB returns the singleton Postgres client
func GetDB(cfg config.Config) (*DB, error) {
	var err error
	once.Do(func() {
		instance, err = connectToDB(cfg)
	})
	return instance, err
}

// ConnectToDB initializes a Postgres connection pool and sets up graceful shutdown.
func connectToDB(cfg config.Config) (*DB, error) {
	connStr, err := buildDSN(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build DSN: %w", err)
	}

	db, err := sqlx.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	setConnectionPool(db)

	wrappedDB := &DB{DB: db}
	return wrappedDB, nil
}

// buildDSN creates a DSN string and registers pgx config.
func buildDSN(cfg config.Config) (string, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DB,
	)

	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return "", fmt.Errorf("invalid DSN: %w", err)
	}

	return stdlib.RegisterConnConfig(connConfig), nil
}

// setConnectionPool sets recommended connection pool settings.
func setConnectionPool(db *sqlx.DB) {
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Minute)
	db.SetConnMaxLifetime(5 * time.Minute)
}
