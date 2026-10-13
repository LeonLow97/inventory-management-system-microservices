package usecases

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/observability/logging"
)

type userOutboundStub struct {
	loginUser     *domain.User
	loginToken    string
	loginErr      error
	signUpErr     error
	users         []domain.User
	nextCursor    string
	getUsersErr   error
	updateUserErr error
	loginInput    domain.User
	signUpInput   domain.User
	usersLimit    int64
	usersCursor   string
	updateInput   domain.User
}

type loginAttemptStub struct{}

func (loginAttemptStub) LoginRateLimiter(context.Context, string, string) error { return nil }
func (loginAttemptStub) ResetLoginFailRateLimiter(context.Context, string, string) error {
	return nil
}
func (loginAttemptStub) IncrementLoginFailRateLimiter(context.Context, string, string) error {
	return nil
}

func (s *userOutboundStub) Login(_ context.Context, req domain.User) (*domain.User, string, error) {
	s.loginInput = req
	return s.loginUser, s.loginToken, s.loginErr
}

func (s *userOutboundStub) SignUp(_ context.Context, req domain.User) error {
	s.signUpInput = req
	return s.signUpErr
}

func (s *userOutboundStub) GetUsers(_ context.Context, limit int64, cursor string) ([]domain.User, string, error) {
	s.usersLimit = limit
	s.usersCursor = cursor
	return s.users, s.nextCursor, s.getUsersErr
}

func (s *userOutboundStub) UpdateUser(_ context.Context, req domain.User) error {
	s.updateInput = req
	return s.updateUserErr
}

func TestServiceLogin(t *testing.T) {
	cause := errors.New("login failed")
	user := &domain.User{ID: 3, Email: "jane@example.com"}
	request := domain.User{Email: "jane@example.com", Password: "Secret123!"}
	testCases := []struct {
		name      string
		repoUser  *domain.User
		repoToken string
		repoErr   error
		wantUser  *domain.User
		wantToken string
		wantErr   error
	}{
		{name: "success", repoUser: user, repoToken: "jwt-token", wantUser: user, wantToken: "jwt-token"},
		{name: "failure", repoErr: cause, wantErr: cause},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			outbound := &userOutboundStub{loginUser: testCase.repoUser, loginToken: testCase.repoToken, loginErr: testCase.repoErr}
			service := NewUserService(outbound, loginAttemptStub{}, logging.New(io.Discard))

			gotUser, gotToken, gotErr := service.Login(context.Background(), request, "192.0.2.1")
			if !reflect.DeepEqual(gotUser, testCase.wantUser) || gotToken != testCase.wantToken || !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("Login() = (%#v, %q, %v), want (%#v, %q, %v)", gotUser, gotToken, gotErr, testCase.wantUser, testCase.wantToken, testCase.wantErr)
			}
			if !reflect.DeepEqual(outbound.loginInput, request) {
				t.Errorf("outbound received %#v, want %#v", outbound.loginInput, request)
			}
		})
	}
}

func TestServiceSignUp(t *testing.T) {
	cause := errors.New("signup failed")
	request := domain.User{Email: "jane@example.com", Password: "Secret123!"}
	testCases := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "success"},
		{name: "failure", repoErr: cause, wantErr: cause},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			outbound := &userOutboundStub{signUpErr: testCase.repoErr}
			service := NewUserService(outbound, loginAttemptStub{}, logging.New(io.Discard))
			gotErr := service.SignUp(context.Background(), request)
			if !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("SignUp() error = %v, want %v", gotErr, testCase.wantErr)
			}
			if !reflect.DeepEqual(outbound.signUpInput, request) {
				t.Errorf("outbound received %#v, want %#v", outbound.signUpInput, request)
			}
		})
	}
}

func TestServiceGetUsers(t *testing.T) {
	cause := errors.New("list users failed")
	users := []domain.User{{ID: 1}, {ID: 2}}
	testCases := []struct {
		name       string
		repoUsers  []domain.User
		nextCursor string
		repoErr    error
		wantUsers  []domain.User
		wantCursor string
		wantErr    error
	}{
		{name: "success", repoUsers: users, nextCursor: "next", wantUsers: users, wantCursor: "next"},
		{name: "failure", repoErr: cause, wantErr: cause},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			outbound := &userOutboundStub{users: testCase.repoUsers, nextCursor: testCase.nextCursor, getUsersErr: testCase.repoErr}
			service := NewUserService(outbound, loginAttemptStub{}, logging.New(io.Discard))
			gotUsers, gotCursor, gotErr := service.GetUsers(context.Background(), 25, "current")
			if !reflect.DeepEqual(gotUsers, testCase.wantUsers) || gotCursor != testCase.wantCursor || !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("GetUsers() = (%#v, %q, %v), want (%#v, %q, %v)", gotUsers, gotCursor, gotErr, testCase.wantUsers, testCase.wantCursor, testCase.wantErr)
			}
			if outbound.usersLimit != 25 || outbound.usersCursor != "current" {
				t.Errorf("outbound received pagination (%d, %q), want (25, %q)", outbound.usersLimit, outbound.usersCursor, "current")
			}
		})
	}
}

func TestServiceUpdateUser(t *testing.T) {
	cause := errors.New("update failed")
	request := domain.User{FirstName: "Jane", LastName: "Doe"}
	testCases := []struct {
		name    string
		repoErr error
		wantErr error
	}{
		{name: "success"},
		{name: "failure", repoErr: cause, wantErr: cause},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			outbound := &userOutboundStub{updateUserErr: testCase.repoErr}
			service := NewUserService(outbound, loginAttemptStub{}, logging.New(io.Discard))
			gotErr := service.UpdateUser(context.Background(), request)
			if !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("UpdateUser() error = %v, want %v", gotErr, testCase.wantErr)
			}
			if !reflect.DeepEqual(outbound.updateInput, request) {
				t.Errorf("outbound received %#v, want %#v", outbound.updateInput, request)
			}
		})
	}
}

func (s *userOutboundStub) GetCurrentUser(context.Context) (*domain.User, error) {
	return s.loginUser, s.loginErr
}
