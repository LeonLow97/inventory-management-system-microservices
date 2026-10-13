package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
)

func TestGetCurrentUser(t *testing.T) {
	dbErr := errors.New("database unavailable")
	cases := []struct {
		name    string
		id      int64
		user    *domain.User
		repoErr error
		wantErr error
		calls   int
	}{
		{"active user", 42, &domain.User{ID: 42, Email: "ada@example.com", Active: true}, nil, nil, 1},
		{"inactive user", 42, &domain.User{ID: 42, Active: false}, nil, coreerrors.ErrInactiveUser, 1},
		{"missing user", 42, nil, coreerrors.ErrUserNotFound, coreerrors.ErrUserNotFound, 1},
		{"database failure", 42, nil, dbErr, dbErr, 1},
		{"invalid identity", 0, nil, nil, coreerrors.ErrInvalidCredentials, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			repo := &userRepositoryStub{getUserByEmailUser: tc.user, getUserByEmailErr: tc.repoErr}
			service := NewService(Config{}, repo, &cacheStub{}, testLogger())
			got, err := service.GetCurrentUser(context.Background(), tc.id)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("error = %v, want %v", err, tc.wantErr)
			}
			if tc.wantErr == nil && got != tc.user {
				t.Fatalf("user = %#v, want database user", got)
			}
			if tc.wantErr != nil && got != nil {
				t.Fatal("must not return a profile on failure")
			}
			if repo.currentUserCalls != tc.calls || (tc.calls > 0 && repo.currentUserID != tc.id) {
				t.Fatal("incorrect database identity lookup")
			}
		})
	}
}
