package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/jmoiron/sqlx"
)

func TestUserRepositoryGetUserByEmail(t *testing.T) {
	query := `SELECT id, email, hashed_password, first_name, last_name, active, created_at, updated_at FROM users WHERE email = $1`
	createdAt := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	databaseErr := errors.New("query failed")
	testCases := []struct {
		name    string
		dbErr   error
		wantErr error
	}{
		{name: "returns user"},
		{name: "maps missing row to application error", dbErr: sql.ErrNoRows, wantErr: coreerrors.ErrUserNotFound},
		{name: "propagates database error", dbErr: databaseErr, wantErr: databaseErr},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository, mock := newUserRepositoryMock(t)
			expectation := mock.ExpectQuery(query).WithArgs("ada@example.com")
			if testCase.dbErr != nil {
				expectation.WillReturnError(testCase.dbErr)
			} else {
				expectation.WillReturnRows(sqlmock.NewRows([]string{
					"id", "email", "hashed_password", "first_name", "last_name", "active", "created_at", "updated_at",
				}).AddRow(42, "ada@example.com", "hashed", "Ada", "Lovelace", true, createdAt, updatedAt))
			}

			got, err := repository.GetUserByEmail(context.Background(), "ada@example.com")
			if testCase.wantErr != nil {
				if !errors.Is(err, testCase.wantErr) {
					t.Fatalf("GetUserByEmail() error = %v, want %v", err, testCase.wantErr)
				}
				if got != nil {
					t.Errorf("GetUserByEmail() user = %#v, want nil on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetUserByEmail() error = %v", err)
			}
			if got == nil {
				t.Fatal("GetUserByEmail() user = nil, want mapped user")
			}
			if got.ID != 42 || got.Email != "ada@example.com" || got.HashedPassword == nil || *got.HashedPassword != "hashed" || got.FirstName == nil || *got.FirstName != "Ada" || got.LastName == nil || *got.LastName != "Lovelace" || !got.Active || !got.CreatedAt.Equal(createdAt) || !got.UpdatedAt.Equal(updatedAt) {
				t.Errorf("GetUserByEmail() = %#v, want database row fields mapped", got)
			}
		})
	}
}

func TestUserRepositoryEmailExists(t *testing.T) {
	testCases := []struct {
		name    string
		exists  bool
		dbErr   error
		wantErr bool
	}{
		{name: "email exists", exists: true},
		{name: "email does not exist"},
		{name: "query error", dbErr: errors.New("query failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository, mock := newUserRepositoryMock(t)
			expectation := mock.ExpectQuery(`SELECT EXISTS ( SELECT 1 FROM users WHERE email = $1 )`).WithArgs("ada@example.com")
			if testCase.dbErr != nil {
				expectation.WillReturnError(testCase.dbErr)
			} else {
				expectation.WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(testCase.exists))
			}

			got, err := repository.EmailExists(context.Background(), "ada@example.com")
			if got != testCase.exists || (err != nil) != testCase.wantErr {
				t.Errorf("EmailExists() = (%t, %v), want (%t, err=%t)", got, err, testCase.exists, testCase.wantErr)
			}
		})
	}
}

func TestUserRepositoryInsertUser(t *testing.T) {
	user := &domain.User{
		Email:          "ada@example.com",
		HashedPassword: stringPointerForRepositoryTest("hashed"),
		FirstName:      stringPointerForRepositoryTest("Ada"),
		LastName:       stringPointerForRepositoryTest("Lovelace"),
	}
	testCases := []struct {
		name    string
		dbErr   error
		wantErr bool
	}{
		{name: "inserts user"},
		{name: "returns insert error", dbErr: errors.New("insert failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository, mock := newUserRepositoryMock(t)
			expectation := mock.ExpectExec(`INSERT INTO users (email, hashed_password, first_name, last_name) VALUES ($1, $2, $3, $4)`).
				WithArgs(user.Email, *user.HashedPassword, *user.FirstName, *user.LastName)
			if testCase.dbErr != nil {
				expectation.WillReturnError(testCase.dbErr)
			} else {
				expectation.WillReturnResult(sqlmock.NewResult(1, 1))
			}

			err := repository.InsertUser(context.Background(), user)
			if (err != nil) != testCase.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %t", err, testCase.wantErr)
			}
		})
	}
}

func TestUserRepositoryUpdateUserByEmail(t *testing.T) {
	updateErr := errors.New("update failed")
	testCases := []struct {
		name     string
		user     *domain.User
		query    string
		args     []driver.Value
		dbErr    error
		wantErr  bool
		wantExec bool
	}{
		{name: "requires email", user: &domain.User{}, wantErr: true},
		{name: "no fields is a no-op", user: &domain.User{Email: "ada@example.com"}},
		{
			name: "updates provided fields",
			user: &domain.User{
				Email: "ada@example.com", FirstName: stringPointerForRepositoryTest("Ada"),
				LastName: stringPointerForRepositoryTest("Lovelace"), HashedPassword: stringPointerForRepositoryTest("hashed"),
			},
			query:    `UPDATE users SET first_name = $1, last_name = $2, hashed_password = $3, updated_at = NOW() WHERE email = $4`,
			args:     []driver.Value{"Ada", "Lovelace", "hashed", "ada@example.com"},
			wantExec: true,
		},
		{
			name:  "propagates update error",
			user:  &domain.User{Email: "ada@example.com", FirstName: stringPointerForRepositoryTest("Ada")},
			query: `UPDATE users SET first_name = $1, updated_at = NOW() WHERE email = $2`,
			args:  []driver.Value{"Ada", "ada@example.com"}, dbErr: updateErr, wantErr: true, wantExec: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository, mock := newUserRepositoryMock(t)
			if testCase.wantExec {
				expectation := mock.ExpectExec(testCase.query).WithArgs(testCase.args...)
				if testCase.dbErr != nil {
					expectation.WillReturnError(testCase.dbErr)
				} else {
					expectation.WillReturnResult(sqlmock.NewResult(0, 1))
				}
			}

			err := repository.UpdateUserByEmail(context.Background(), testCase.user)
			if (err != nil) != testCase.wantErr {
				t.Errorf("UpdateUserByEmail() error = %v, wantErr %t", err, testCase.wantErr)
			}
		})
	}
}

func TestUserRepositoryGetUsers(t *testing.T) {
	testCases := []struct {
		name    string
		cursor  domain.UserCursor
		query   string
		args    []driver.Value
		dbErr   error
		wantErr bool
		wantRow bool
	}{
		{
			name:  "first page",
			query: `SELECT id, email, first_name, last_name, active, updated_at, created_at FROM users ORDER BY id DESC LIMIT $1`,
			args:  []driver.Value{int64(10)}, wantRow: true,
		},
		{
			name:   "subsequent page",
			cursor: domain.UserCursor{ID: 20},
			query:  `SELECT id, email, first_name, last_name, active, updated_at, created_at FROM users WHERE id < $1 ORDER BY id DESC LIMIT $2`,
			args:   []driver.Value{int64(20), int64(10)}, wantRow: true,
		},
		{
			name:  "empty page",
			query: `SELECT id, email, first_name, last_name, active, updated_at, created_at FROM users ORDER BY id DESC LIMIT $1`,
			args:  []driver.Value{int64(10)},
		},
		{
			name:  "query error",
			query: `SELECT id, email, first_name, last_name, active, updated_at, created_at FROM users ORDER BY id DESC LIMIT $1`,
			args:  []driver.Value{int64(10)}, dbErr: errors.New("query failed"), wantErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository, mock := newUserRepositoryMock(t)
			expectation := mock.ExpectQuery(testCase.query).WithArgs(testCase.args...)
			if testCase.dbErr != nil {
				expectation.WillReturnError(testCase.dbErr)
			} else if testCase.wantRow {
				expectation.WillReturnRows(sqlmock.NewRows([]string{
					"id", "email", "first_name", "last_name", "active", "updated_at", "created_at",
				}).AddRow(12, "ada@example.com", "Ada", "Lovelace", true, time.Time{}, time.Time{}))
			} else {
				expectation.WillReturnRows(sqlmock.NewRows([]string{"id", "email", "first_name", "last_name", "active", "updated_at", "created_at"}))
			}

			got, err := repository.GetUsers(context.Background(), 10, testCase.cursor)
			if (err != nil) != testCase.wantErr {
				t.Fatalf("GetUsers() error = %v, wantErr %t", err, testCase.wantErr)
			}
			if testCase.wantRow && (len(got) != 1 || got[0].ID != 12 || got[0].Email != "ada@example.com") {
				t.Errorf("GetUsers() = %#v, want one mapped user", got)
			}
		})
	}
}

func TestUserRepositoryIsAdminUser(t *testing.T) {
	testCases := []struct {
		name    string
		isAdmin bool
		dbErr   error
		wantErr bool
	}{
		{name: "admin user", isAdmin: true},
		{name: "non-admin user"},
		{name: "query error", dbErr: errors.New("query failed"), wantErr: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository, mock := newUserRepositoryMock(t)
			expectation := mock.ExpectQuery(`SELECT EXISTS ( SELECT 1 FROM admin_users WHERE active AND user_id = $1 )`).WithArgs(int64(7))
			if testCase.dbErr != nil {
				expectation.WillReturnError(testCase.dbErr)
			} else {
				expectation.WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(testCase.isAdmin))
			}

			got, err := repository.IsAdminUser(context.Background(), 7)
			if got != testCase.isAdmin || (err != nil) != testCase.wantErr {
				t.Errorf("IsAdminUser() = (%t, %v), want (%t, err=%t)", got, err, testCase.isAdmin, testCase.wantErr)
			}
		})
	}
}

func newUserRepositoryMock(t *testing.T) (*UserRepository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("sqlmock.New() error = %v", err)
	}
	sqlxDB := sqlx.NewDb(db, "postgres")
	mock.MatchExpectationsInOrder(false)
	mock.ExpectClose()
	t.Cleanup(func() {
		if err := sqlxDB.Close(); err != nil {
			t.Errorf("close mock database: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet SQL expectations: %v", err)
		}
	})

	return &UserRepository{db: sqlxDB}, mock
}

func stringPointerForRepositoryTest(value string) *string {
	return &value
}
