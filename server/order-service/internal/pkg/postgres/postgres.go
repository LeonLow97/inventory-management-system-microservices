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
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresConfig.User,
		cfg.PostgresConfig.Password,
		cfg.PostgresConfig.Host,
		cfg.PostgresConfig.Port,
		cfg.PostgresConfig.DB,
	)

	connConfig, err := pgx.ParseConfig(databaseURL)
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

	return db
}
