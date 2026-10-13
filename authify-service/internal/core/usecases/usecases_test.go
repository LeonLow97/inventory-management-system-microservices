package usecases

import (
	"context"
	"io"

	"github.com/LeonLow97/internal/core/domain"
	"github.com/LeonLow97/internal/core/ports"
	"github.com/LeonLow97/internal/observability/logging"
)

func testLogger() ports.Logger {
	return logging.New(io.Discard)
}

type userRepositoryStub struct {
	currentUserID      int64
	currentUserCalls   int
	getUserByEmailUser *domain.User
	getUserByEmailErr  error
	getUsersUsers      []domain.User
	getUsersErr        error
	getUsersCalls      int
	getUsersLimit      int64
	getUsersCursor     domain.UserCursor
	emailExists        bool
	emailExistsErr     error
	insertUser         *domain.User
	insertUserErr      error
	updateUser         *domain.User
	updateUserErr      error
	isAdmin            bool
	isAdminErr         error
}

func (s *userRepositoryStub) GetUserByEmail(context.Context, string) (*domain.User, error) {
	return s.getUserByEmailUser, s.getUserByEmailErr
}

func (s *userRepositoryStub) GetUsers(_ context.Context, limit int64, cursor domain.UserCursor) ([]domain.User, error) {
	s.getUsersCalls++
	s.getUsersLimit = limit
	s.getUsersCursor = cursor
	return s.getUsersUsers, s.getUsersErr
}

func (s *userRepositoryStub) EmailExists(context.Context, string) (bool, error) {
	return s.emailExists, s.emailExistsErr
}

func (s *userRepositoryStub) InsertUser(_ context.Context, user *domain.User) error {
	s.insertUser = user
	return s.insertUserErr
}

func (s *userRepositoryStub) UpdateUserByEmail(_ context.Context, user *domain.User) error {
	s.updateUser = user
	return s.updateUserErr
}

func (s *userRepositoryStub) IsAdminUser(context.Context, int64) (bool, error) {
	return s.isAdmin, s.isAdminErr
}

type cacheStub struct {
	setNXCalls  int
	setNXKey    string
	setNXValue  any
	setNXTTL    int64
	setNXResult bool
	setNXErr    error
}

func (*cacheStub) Set(context.Context, string, any, int64) error { return nil }

func (s *cacheStub) SetNX(_ context.Context, key string, value any, ttlSeconds int64) (bool, error) {
	s.setNXCalls++
	s.setNXKey = key
	s.setNXValue = value
	s.setNXTTL = ttlSeconds
	return s.setNXResult, s.setNXErr
}

func (*cacheStub) Get(context.Context, string) (string, error) { return "", nil }

func (*cacheStub) Delete(context.Context, string) error { return nil }

func (*cacheStub) Exists(context.Context, string) (bool, error) { return false, nil }

func (*cacheStub) Ping(context.Context) error { return nil }

func (s *userRepositoryStub) GetUserByID(_ context.Context, userID int64) (*domain.User, error) {
	s.currentUserID = userID
	s.currentUserCalls++
	return s.getUserByEmailUser, s.getUserByEmailErr
}
