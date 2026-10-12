package usecases

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/LeonLow97/internal/pkg/cursor"
	"golang.org/x/crypto/bcrypt"
)

func TestServiceGetUsers(t *testing.T) {
	adminErr := errors.New("admin lookup failed")
	queryErr := errors.New("user query failed")
	encodedCursor, err := cursor.EncodeCursor(domain.UserCursor{ID: 25})
	if err != nil {
		t.Fatalf("EncodeCursor() error = %v", err)
	}
	users := []domain.User{{ID: 11}, {ID: 7}}

	testCases := []struct {
		name              string
		isAdmin           bool
		isAdminErr        error
		cursor            string
		users             []domain.User
		getUsersErr       error
		wantErr           error
		wantErrorText     string
		wantGetUsersCalls int
		wantCursorID      int64
		wantNextCursor    bool
	}{
		{name: "returns users and encodes next cursor", isAdmin: true, users: users, wantGetUsersCalls: 1, wantCursorID: 7, wantNextCursor: true},
		{name: "passes decoded cursor to repository", isAdmin: true, cursor: encodedCursor, wantGetUsersCalls: 1, wantNextCursor: false},
		{name: "returns empty page without next cursor", isAdmin: true, wantGetUsersCalls: 1},
		{name: "propagates admin lookup error", isAdminErr: adminErr, wantErr: adminErr},
		{name: "rejects non-admin user", wantErr: coreerrors.ErrUnauthorizedAdminAccess},
		{name: "rejects malformed cursor", isAdmin: true, cursor: "not-base64", wantErrorText: "failed to decode cursor", wantGetUsersCalls: 0},
		{name: "propagates repository query error", isAdmin: true, getUsersErr: queryErr, wantErr: queryErr, wantGetUsersCalls: 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository := &userRepositoryStub{
				isAdmin:       testCase.isAdmin,
				isAdminErr:    testCase.isAdminErr,
				getUsersUsers: testCase.users,
				getUsersErr:   testCase.getUsersErr,
			}
			service := NewService(Config{}, repository, &cacheStub{}, testLogger())

			gotUsers, gotCursor, gotErr := service.GetUsers(context.Background(), 99, 20, testCase.cursor)
			if testCase.wantErr != nil && !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("GetUsers() error = %v, want error matching %v", gotErr, testCase.wantErr)
			}
			if testCase.wantErrorText != "" && (gotErr == nil || !strings.Contains(gotErr.Error(), testCase.wantErrorText)) {
				t.Fatalf("GetUsers() error = %v, want message containing %q", gotErr, testCase.wantErrorText)
			}
			if testCase.wantErr == nil && testCase.wantErrorText == "" && gotErr != nil {
				t.Fatalf("GetUsers() error = %v, want nil", gotErr)
			}
			if repository.getUsersCalls != testCase.wantGetUsersCalls {
				t.Fatalf("repository GetUsers calls = %d, want %d", repository.getUsersCalls, testCase.wantGetUsersCalls)
			}
			if testCase.wantGetUsersCalls == 0 {
				if gotUsers != nil || gotCursor != "" {
					t.Errorf("GetUsers() = (%#v, %q), want (nil, empty cursor) on error", gotUsers, gotCursor)
				}
				return
			}
			if repository.getUsersLimit != 20 {
				t.Errorf("repository limit = %d, want 20", repository.getUsersLimit)
			}
			if testCase.cursor != "" && repository.getUsersCursor.ID != 25 {
				t.Errorf("repository cursor ID = %d, want 25", repository.getUsersCursor.ID)
			}
			if testCase.wantErr != nil || testCase.wantErrorText != "" {
				return
			}
			if !reflect.DeepEqual(gotUsers, testCase.users) {
				t.Errorf("GetUsers() users = %#v, want %#v", gotUsers, testCase.users)
			}
			if !testCase.wantNextCursor {
				if gotCursor != "" {
					t.Errorf("GetUsers() cursor = %q, want empty", gotCursor)
				}
				return
			}
			var nextCursor domain.UserCursor
			if err := cursor.DecodeCursor(gotCursor, &nextCursor); err != nil {
				t.Fatalf("DecodeCursor() error = %v", err)
			}
			if nextCursor.ID != testCase.wantCursorID {
				t.Errorf("next cursor ID = %d, want %d", nextCursor.ID, testCase.wantCursorID)
			}
		})
	}
}

func TestServiceUpdateUserByEmail(t *testing.T) {
	updateErr := errors.New("update failed")
	passwordTooLong := strings.Repeat("x", 73)

	testCases := []struct {
		name               string
		input              domain.UpdateUserInput
		repositoryErr      error
		wantErr            error
		wantRepositoryCall bool
		wantPasswordHash   bool
	}{
		{name: "updates profile fields", input: domain.UpdateUserInput{FirstName: stringPointer("Ada"), LastName: stringPointer("Lovelace")}, wantRepositoryCall: true},
		{name: "hashes a new password", input: domain.UpdateUserInput{Password: stringPointer("NewSecret123!")}, wantRepositoryCall: true, wantPasswordHash: true},
		{name: "propagates repository error", input: domain.UpdateUserInput{FirstName: stringPointer("Ada")}, repositoryErr: updateErr, wantErr: updateErr, wantRepositoryCall: true},
		{name: "rejects password longer than bcrypt limit", input: domain.UpdateUserInput{Password: &passwordTooLong}, wantErr: bcrypt.ErrPasswordTooLong},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository := &userRepositoryStub{updateUserErr: testCase.repositoryErr}
			service := NewService(Config{}, repository, &cacheStub{}, testLogger())
			gotErr := service.UpdateUserByEmail(context.Background(), 42, testCase.input)

			if testCase.wantErr != nil && !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("UpdateUserByEmail() error = %v, want error matching %v", gotErr, testCase.wantErr)
			}
			if testCase.wantErr == nil && gotErr != nil {
				t.Fatalf("UpdateUserByEmail() error = %v, want nil", gotErr)
			}
			if (repository.updateUser != nil) != testCase.wantRepositoryCall {
				t.Fatalf("repository update user = %#v, want call = %t", repository.updateUser, testCase.wantRepositoryCall)
			}
			if !testCase.wantRepositoryCall {
				return
			}
			if repository.updateUser.ID != 42 || !reflect.DeepEqual(repository.updateUser.FirstName, testCase.input.FirstName) || !reflect.DeepEqual(repository.updateUser.LastName, testCase.input.LastName) {
				t.Errorf("repository update user = %#v, want ID and profile fields copied", repository.updateUser)
			}
			if testCase.wantPasswordHash {
				if repository.updateUser.HashedPassword == nil {
					t.Fatal("repository update user has no password hash")
				}
				if err := bcrypt.CompareHashAndPassword([]byte(*repository.updateUser.HashedPassword), []byte(*testCase.input.Password)); err != nil {
					t.Errorf("stored password hash does not match input: %v", err)
				}
			} else if repository.updateUser.HashedPassword != nil {
				t.Errorf("repository password hash = %q, want nil", *repository.updateUser.HashedPassword)
			}
		})
	}
}
