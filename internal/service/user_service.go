package service

import (
	"database/sql"
	"errors"
	"wallets-api-postgres/internal/models"
	"wallets-api-postgres/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) CreateUser(input models.RegisterInput) (*models.User, error) {
	if input.Email == "" {
		return nil, ErrEmailRequired
	}

	if input.Password == "" {
		return nil, ErrPasswordRequired
	}

	if len(input.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	_, err := s.userRepository.GetUserByEmail(input.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hashing password!")
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	createdUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}
