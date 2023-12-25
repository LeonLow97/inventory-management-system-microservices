package users

import (
	"github.com/LeonLow97/utils"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	UpdateUser(req UpdateUserRequestDTO) error
	GetUsers() (*[]User, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{
		repo: r,
	}
}

func (s service) UpdateUser(req UpdateUserRequestDTO) error {
	// check if username exists
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return err
	}

	// check if FirstName, LastName, Password, Email are the same as previous
	if req.FirstName == user.FirstName || req.LastName == user.LastName || req.Email == user.Email {
		return ErrSameValue
	}

	// check if password is the same as previous
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		if err != bcrypt.ErrMismatchedHashAndPassword {
			return err
		}
	} else {
		return ErrSameValue
	}

	// check if password contains at least 1 uppercase, lowercase, numeric and special character
	if !utils.IsValidPassword(req.Password) {
		return ErrInvalidPasswordFormat
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	req.Password = string(hashedPassword)

	// update user in database, set updated_at time to now
	if err = s.repo.UpdateUserByUsername(req); err != nil {
		return err
	}

	return nil
}

func (s service) GetUsers() (*[]User, error) {
	// get all users from the database
	users, err := s.repo.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}
