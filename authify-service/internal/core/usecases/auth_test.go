package usecases

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/LeonLow97/internal/core/domain"
	coreerrors "github.com/LeonLow97/internal/core/errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func TestServiceLogin(t *testing.T) {
	const password = "Secret123!"
	passwordHash := hashPasswordForTest(t, password)
	secret := "test-signing-secret"
	validUser := &domain.User{
		ID:             42,
		Email:          "ada@example.com",
		HashedPassword: &passwordHash,
		Active:         true,
	}
	cacheErr := errors.New("cache unavailable")
	repositoryErr := errors.New("database unavailable")

	testCases := []struct {
		name          string
		user          *domain.User
		password      string
		repositoryErr error
		cacheErr      error
		wantErr       error
		wantErrorText string
		wantCacheCall int
		wantSuccess   bool
	}{
		{name: "creates session and token", user: validUser, password: password, wantCacheCall: 1, wantSuccess: true},
		{name: "invalid password", user: validUser, password: "incorrect", wantErr: coreerrors.ErrInvalidCredentials},
		{name: "invalid password hash", user: &domain.User{ID: 42, HashedPassword: stringPointer("bad"), Active: true}, password: password, wantErr: coreerrors.ErrInvalidCredentials},
		{name: "inactive user", user: &domain.User{ID: 42, HashedPassword: &passwordHash}, password: password, wantErr: coreerrors.ErrInactiveUser},
		{name: "repository failure", repositoryErr: repositoryErr, wantErrorText: "look up user for login"},
		{name: "cache failure", user: validUser, password: password, cacheErr: cacheErr, wantErr: cacheErr, wantCacheCall: 1},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository := &userRepositoryStub{
				getUserByEmailUser: testCase.user,
				getUserByEmailErr:  testCase.repositoryErr,
			}
			appCache := &cacheStub{setNXErr: testCase.cacheErr, setNXResult: true}
			service := NewService(Config{JWTConfig: JWTConfig{SecretKey: secret, Expiry: 15}}, repository, appCache, testLogger())

			gotUser, gotToken, gotErr := service.Login(context.Background(), domain.LoginInput{
				Email:    "ada@example.com",
				Password: testCase.password,
			})

			if testCase.wantErr != nil && !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("Login() error = %v, want error matching %v", gotErr, testCase.wantErr)
			}
			if testCase.wantErrorText != "" && (gotErr == nil || !strings.Contains(gotErr.Error(), testCase.wantErrorText)) {
				t.Fatalf("Login() error = %v, want message containing %q", gotErr, testCase.wantErrorText)
			}
			if testCase.wantErr == nil && testCase.wantErrorText == "" && gotErr != nil {
				t.Fatalf("Login() error = %v, want nil", gotErr)
			}
			if appCache.setNXCalls != testCase.wantCacheCall {
				t.Fatalf("cache SetNX calls = %d, want %d", appCache.setNXCalls, testCase.wantCacheCall)
			}

			if !testCase.wantSuccess {
				if gotUser != nil || gotToken != "" {
					t.Fatalf("Login() = (%#v, %q), want (nil, empty token) on error", gotUser, gotToken)
				}
				return
			}
			if gotUser != validUser || gotToken == "" {
				t.Fatalf("Login() = (%#v, %q), want user and signed token", gotUser, gotToken)
			}
			if appCache.setNXKey == "" || !strings.HasPrefix(appCache.setNXKey, "user_session:") {
				t.Errorf("session cache key = %q, want user_session prefix", appCache.setNXKey)
			}
			if appCache.setNXValue != validUser.ID || appCache.setNXTTL != int64(16*time.Minute/time.Second) {
				t.Errorf("session cache values = (%v, %d), want (%d, %d)", appCache.setNXValue, appCache.setNXTTL, validUser.ID, int64(16*time.Minute/time.Second))
			}

			claims := &CustomClaims{}
			parsedToken, err := jwt.ParseWithClaims(gotToken, claims, func(*jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err != nil || !parsedToken.Valid {
				t.Fatalf("ParseWithClaims() error = %v, valid = %t", err, parsedToken != nil && parsedToken.Valid)
			}
			if claims.UserID != validUser.ID || claims.SessionID == "" || claims.Issuer != "authify-service" {
				t.Errorf("token claims = %#v, want user id %d, session id, and issuer", claims, validUser.ID)
			}
		})
	}
}

func TestServiceSignUp(t *testing.T) {
	databaseErr := errors.New("insert failed")
	lookupErr := errors.New("lookup failed")
	testCases := []struct {
		name           string
		password       string
		emailExists    bool
		emailErr       error
		insertErr      error
		wantErr        error
		wantErrorText  string
		wantInsertCall bool
	}{
		{name: "creates user", password: "Secret123!", wantInsertCall: true},
		{name: "rejects duplicate email", emailExists: true, wantErr: coreerrors.ErrEmailAlreadyExists},
		{name: "email lookup failure", emailErr: lookupErr, wantErrorText: "check whether signup email exists"},
		{name: "insert failure", password: "Secret123!", insertErr: databaseErr, wantErrorText: "create user account", wantInsertCall: true},
		{name: "rejects password longer than bcrypt limit", password: strings.Repeat("x", 73), wantErrorText: "hash signup password"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repository := &userRepositoryStub{
				emailExists:    testCase.emailExists,
				emailExistsErr: testCase.emailErr,
				insertUserErr:  testCase.insertErr,
			}
			service := NewService(Config{}, repository, &cacheStub{}, testLogger())
			gotErr := service.SignUp(context.Background(), domain.SignUpInput{
				Email:     "ada@example.com",
				Password:  testCase.password,
				FirstName: "Ada",
				LastName:  "Lovelace",
			})

			if testCase.wantErr != nil && !errors.Is(gotErr, testCase.wantErr) {
				t.Fatalf("SignUp() error = %v, want error matching %v", gotErr, testCase.wantErr)
			}
			if testCase.wantErrorText != "" && (gotErr == nil || !strings.Contains(gotErr.Error(), testCase.wantErrorText)) {
				t.Fatalf("SignUp() error = %v, want message containing %q", gotErr, testCase.wantErrorText)
			}
			if testCase.wantErr == nil && testCase.wantErrorText == "" && gotErr != nil {
				t.Fatalf("SignUp() error = %v, want nil", gotErr)
			}
			if (repository.insertUser != nil) != testCase.wantInsertCall {
				t.Fatalf("inserted user = %#v, want insert call = %t", repository.insertUser, testCase.wantInsertCall)
			}
			if !testCase.wantInsertCall {
				return
			}
			insertedUser := repository.insertUser
			if insertedUser == nil || insertedUser.HashedPassword == nil || insertedUser.FirstName == nil || insertedUser.LastName == nil {
				t.Fatalf("inserted user = %#v, want complete user fields", insertedUser)
			}
			if insertedUser.Email != "ada@example.com" || *insertedUser.FirstName != "Ada" || *insertedUser.LastName != "Lovelace" {
				t.Errorf("inserted user = %#v, want signup fields copied", insertedUser)
			}
			if err := bcrypt.CompareHashAndPassword([]byte(*insertedUser.HashedPassword), []byte(testCase.password)); err != nil {
				t.Errorf("stored password hash does not match signup password: %v", err)
			}
		})
	}
}

func hashPasswordForTest(t *testing.T, password string) string {
	t.Helper()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}
	return string(hash)
}

func stringPointer(value string) *string {
	return &value
}
