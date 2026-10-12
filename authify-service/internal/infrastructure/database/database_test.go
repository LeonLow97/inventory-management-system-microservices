package database

import (
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LeonLow97/internal/infrastructure/config"
	"github.com/jmoiron/sqlx"
)

func TestBuildDSN(t *testing.T) {
	testCases := []struct {
		name      string
		cfg       config.Config
		wantError string
	}{
		{
			name: "valid postgres config",
			cfg: config.Config{PostgresConfig: config.PostgresConfig{
				User: "service", Password: "secret", Host: "localhost", Port: 5432, DB: "users",
			}},
		},
		{
			name: "invalid port",
			cfg: config.Config{PostgresConfig: config.PostgresConfig{
				User: "service", Password: "secret", Host: "localhost", Port: 70000, DB: "users",
			}},
			wantError: "invalid DSN",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			got, err := buildDSN(testCase.cfg)
			if testCase.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), testCase.wantError) {
					t.Fatalf("buildDSN() error = %v, want message containing %q", err, testCase.wantError)
				}
				if got != "" {
					t.Errorf("buildDSN() = %q, want empty DSN on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("buildDSN() error = %v", err)
			}
			if got == "" {
				t.Error("buildDSN() returned an empty DSN")
			}
		})
	}
}

func TestSetConnectionPool(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	database := sqlx.NewDb(db, "postgres")
	mock.ExpectClose()
	t.Cleanup(func() {
		if err := database.Close(); err != nil {
			t.Errorf("close mock database: %v", err)
		}
	})

	setConnectionPool(database)

	if got := database.Stats().MaxOpenConnections; got != 25 {
		t.Errorf("MaxOpenConnections = %d, want 25", got)
	}
}
