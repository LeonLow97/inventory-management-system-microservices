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

	log.Println("Successfully connected to Postgres database!")
	return db, nil
}
