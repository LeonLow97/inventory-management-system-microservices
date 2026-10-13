package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	coreerrors "github.com/LeonLow97/internal/core/errors"
)

func TestGetUserByID(t *testing.T) {
	dbErr := errors.New("query failed")
	for _, tc := range []struct {
		name           string
		dbErr, wantErr error
	}{
		{"database profile", nil, nil}, {"missing user", sql.ErrNoRows, coreerrors.ErrUserNotFound}, {"database error", dbErr, dbErr},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo, mock := newUserRepositoryMock(t)
			query := mock.ExpectQuery(`SELECT id, email, first_name, last_name, active, created_at, updated_at FROM users WHERE id = $1`).WithArgs(int64(42))
			if tc.dbErr != nil {
				query.WillReturnError(tc.dbErr)
			} else {
				query.WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "active", "created_at", "updated_at"}).AddRow(42, "ada@example.com", "Ada", "Lovelace", true, time.Now(), time.Now()))
			}
			user, err := repo.GetUserByID(context.Background(), 42)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && (user == nil || user.ID != 42 || user.Email != "ada@example.com" || *user.FirstName != "Ada" || *user.LastName != "Lovelace" || user.HashedPassword != nil) {
				t.Fatal("profile must contain database names and exclude password hashes")
			}
			if tc.wantErr != nil && user != nil {
				t.Fatal("profile must be nil on error")
			}
		})
	}
}
