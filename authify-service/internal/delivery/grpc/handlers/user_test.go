package handlers

import (
	"context"
	"errors"
	"testing"

	pb "github.com/LeonLow97/authify/grpc-proto/pb"
	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type userServiceStub struct {
	loginUser   *domain.User
	loginToken  string
	loginErr    error
	loginInput  domain.LoginInput
	loginCalls  int
	signUpErr   error
	signUpInput domain.SignUpInput
	signUpCalls int
	users       []domain.User
	nextCursor  string
	getUsersErr error
	getUsersReq struct {
		adminUserID int64
		limit       int64
		cursor      string
	}
	getUsersCalls int
	updateErr     error
	updateReq     struct {
		userID int64
		input  domain.UpdateUserInput
	}
	updateCalls int
}

func (s *userServiceStub) Login(_ context.Context, input domain.LoginInput) (*domain.User, string, error) {
	s.loginCalls++
	s.loginInput = input
	return s.loginUser, s.loginToken, s.loginErr
}

func (s *userServiceStub) SignUp(_ context.Context, input domain.SignUpInput) error {
	s.signUpCalls++
	s.signUpInput = input
	return s.signUpErr
}

func (s *userServiceStub) UpdateUserByEmail(_ context.Context, userID int64, input domain.UpdateUserInput) error {
	s.updateCalls++
	s.updateReq.userID = userID
	s.updateReq.input = input
	return s.updateErr
}

func (s *userServiceStub) GetUsers(_ context.Context, adminUserID int64, limit int64, cursor string) ([]domain.User, string, error) {
	s.getUsersCalls++
	s.getUsersReq.adminUserID = adminUserID
	s.getUsersReq.limit = limit
	s.getUsersReq.cursor = cursor
	return s.users, s.nextCursor, s.getUsersErr
}

func TestUserHandlerLogin(t *testing.T) {
	user := &domain.User{
		Email:     "ada@example.com",
		FirstName: stringPointerForHandlerTest("Ada"),
		LastName:  stringPointerForHandlerTest("Lovelace"),
		Active:    true,
		Admin:     true,
	}
	serviceErr := errors.New("unexpected failure")
	testCases := []struct {
		name       string
		serviceErr error
		wantCode   codes.Code
		wantCalls  int
	}{
		{name: "success", wantCode: codes.OK, wantCalls: 1},
		{name: "invalid credentials", serviceErr: coreerrors.ErrInvalidCredentials, wantCode: codes.Unauthenticated, wantCalls: 1},
		{name: "unexpected service error", serviceErr: serviceErr, wantCode: codes.Internal, wantCalls: 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{loginUser: user, loginToken: "signed-token", loginErr: testCase.serviceErr}
			handler := NewUserHandler(service)
			response, err := handler.Login(context.Background(), &pb.LoginRequest{Email: " ada@example.com ", Password: " Secret123! "})
			if status.Code(err) != testCase.wantCode {
				t.Fatalf("Login() status = %s, error = %v, want %s", status.Code(err), err, testCase.wantCode)
			}
			if service.loginCalls != testCase.wantCalls {
				t.Fatalf("service Login calls = %d, want %d", service.loginCalls, testCase.wantCalls)
			}
			if service.loginInput != (domain.LoginInput{Email: "ada@example.com", Password: "Secret123!"}) {
				t.Errorf("service login input = %#v, want trimmed request values", service.loginInput)
			}
			if testCase.wantCode == codes.OK && (response == nil || response.Email != user.Email || response.Token != "signed-token") {
				t.Errorf("Login() response = %#v, want mapped user and token", response)
			}
			if testCase.wantCode != codes.OK && response != nil {
				t.Errorf("Login() response = %#v, want nil on error", response)
			}
		})
	}
}

func TestUserHandlerSignUp(t *testing.T) {
	duplicateErr := errors.New("database duplicate")
	testCases := []struct {
		name       string
		serviceErr error
		wantCode   codes.Code
		wantCalls  int
	}{
		{name: "success", wantCode: codes.OK, wantCalls: 1},
		{name: "duplicate email", serviceErr: coreerrors.ErrEmailAlreadyExists, wantCode: codes.AlreadyExists, wantCalls: 1},
		{name: "unexpected service error", serviceErr: duplicateErr, wantCode: codes.Internal, wantCalls: 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{signUpErr: testCase.serviceErr}
			handler := NewUserHandler(service)
			response, err := handler.SignUp(context.Background(), &pb.SignUpRequest{
				Email: " ada@example.com ", Password: " Secret123! ", FirstName: " Ada ", LastName: " Lovelace ",
			})
			if status.Code(err) != testCase.wantCode {
				t.Fatalf("SignUp() status = %s, error = %v, want %s", status.Code(err), err, testCase.wantCode)
			}
			if service.signUpCalls != testCase.wantCalls {
				t.Fatalf("service SignUp calls = %d, want %d", service.signUpCalls, testCase.wantCalls)
			}
			if service.signUpInput != (domain.SignUpInput{
				Email: "ada@example.com", Password: "Secret123!", FirstName: "Ada", LastName: "Lovelace",
			}) {
				t.Errorf("service signup input = %#v, want trimmed request values", service.signUpInput)
			}
			if testCase.wantCode == codes.OK && response == nil {
				t.Error("SignUp() response = nil, want empty response")
			}
			if testCase.wantCode != codes.OK && response != nil {
				t.Errorf("SignUp() response = %#v, want nil on error", response)
			}
		})
	}
}

func TestUserHandlerGetUsers(t *testing.T) {
	serviceErr := errors.New("query failed")
	testCases := []struct {
		name         string
		userID       string
		serviceErr   error
		wantCode     codes.Code
		wantCalls    int
		wantAdminID  int64
		wantLimit    int64
		wantResponse bool
	}{
		{name: "success", userID: "17", wantCode: codes.OK, wantCalls: 1, wantAdminID: 17, wantLimit: 10, wantResponse: true},
		{name: "caps requested limit", userID: "17", wantCode: codes.OK, wantCalls: 1, wantAdminID: 17, wantLimit: 50, wantResponse: true},
		{name: "missing metadata", wantCode: codes.Unauthenticated},
		{name: "admin access denied", userID: "17", serviceErr: coreerrors.ErrUnauthorizedAdminAccess, wantCode: codes.Unauthenticated, wantCalls: 1, wantAdminID: 17, wantLimit: 10},
		{name: "unexpected service error", userID: "17", serviceErr: serviceErr, wantCode: codes.Internal, wantCalls: 1, wantAdminID: 17, wantLimit: 10},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{
				users:       []domain.User{{ID: 8, Email: "ada@example.com", FirstName: stringPointerForHandlerTest("Ada"), LastName: stringPointerForHandlerTest("Lovelace"), Active: true}},
				nextCursor:  "next-page",
				getUsersErr: testCase.serviceErr,
			}
			handler := NewUserHandler(service)
			ctx := context.Background()
			if testCase.userID != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", testCase.userID))
			}
			limit := int64(0)
			if testCase.name == "caps requested limit" {
				limit = 100
			}
			response, err := handler.GetUsers(ctx, &pb.GetUsersRequest{Limit: limit, Cursor: "current"})
			if status.Code(err) != testCase.wantCode {
				t.Fatalf("GetUsers() status = %s, error = %v, want %s", status.Code(err), err, testCase.wantCode)
			}
			if service.getUsersCalls != testCase.wantCalls {
				t.Fatalf("service GetUsers calls = %d, want %d", service.getUsersCalls, testCase.wantCalls)
			}
			if testCase.wantCalls == 1 && (service.getUsersReq.adminUserID != testCase.wantAdminID || service.getUsersReq.limit != testCase.wantLimit || service.getUsersReq.cursor != "current") {
				t.Errorf("service GetUsers arguments = (%d, %d, %q), want (%d, %d, %q)", service.getUsersReq.adminUserID, service.getUsersReq.limit, service.getUsersReq.cursor, testCase.wantAdminID, testCase.wantLimit, "current")
			}
			if testCase.wantResponse && (response == nil || response.NextCursor != "next-page" || len(response.Users) != 1 || response.Users[0].ID != 8) {
				t.Errorf("GetUsers() response = %#v, want mapped users and cursor", response)
			}
			if !testCase.wantResponse && response != nil {
				t.Errorf("GetUsers() response = %#v, want nil on error", response)
			}
		})
	}
}

func TestUserHandlerUpdateUser(t *testing.T) {
	serviceErr := errors.New("update failed")
	testCases := []struct {
		name       string
		userID     string
		serviceErr error
		wantCode   codes.Code
		wantCalls  int
		wantUserID int64
	}{
		{name: "success", userID: "42", wantCode: codes.OK, wantCalls: 1, wantUserID: 42},
		{name: "missing metadata", wantCode: codes.Unauthenticated},
		{name: "user not found", userID: "42", serviceErr: coreerrors.ErrUserNotFound, wantCode: codes.Unauthenticated, wantCalls: 1, wantUserID: 42},
		{name: "unexpected service error", userID: "42", serviceErr: serviceErr, wantCode: codes.Internal, wantCalls: 1, wantUserID: 42},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := &userServiceStub{updateErr: testCase.serviceErr}
			handler := NewUserHandler(service)
			ctx := context.Background()
			if testCase.userID != "" {
				ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", testCase.userID))
			}
			response, err := handler.UpdateUser(ctx, &pb.UpdateUserRequest{FirstName: "Ada", LastName: "Lovelace", Password: "Secret123!"})
			if status.Code(err) != testCase.wantCode {
				t.Fatalf("UpdateUser() status = %s, error = %v, want %s", status.Code(err), err, testCase.wantCode)
			}
			if service.updateCalls != testCase.wantCalls {
				t.Fatalf("service UpdateUserByEmail calls = %d, want %d", service.updateCalls, testCase.wantCalls)
			}
			if testCase.wantCalls == 1 {
				got := service.updateReq
				if got.userID != testCase.wantUserID || got.input.FirstName == nil || *got.input.FirstName != "Ada" || got.input.LastName == nil || *got.input.LastName != "Lovelace" || got.input.Password == nil || *got.input.Password != "Secret123!" {
					t.Errorf("service update request = %#v, want mapped ID and fields", got)
				}
			}
			if testCase.wantCode == codes.OK && response == nil {
				t.Error("UpdateUser() response = nil, want empty response")
			}
			if testCase.wantCode != codes.OK && response != nil {
				t.Errorf("UpdateUser() response = %#v, want nil on error", response)
			}
		})
	}
}
