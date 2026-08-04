package service

import (
	"context"
	"database/sql"
	"errors"
	"wallets-api-postgres/internal/auth"
	"wallets-api-postgres/internal/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, userID int64) (models.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (models.User, error)
}

type UserService struct {
	userRepository UserRepository
	jwtSecret      string
}

func NewUserService(userRepository UserRepository,
	jwtSecret string,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *UserService) CreateUser(ctx context.Context, input models.RegisterInput) (*models.User, error) {
	if input.Email == "" {
		return nil, ErrEmailRequired
	}

	if input.Password == "" {
		return nil, ErrPasswordRequired
	}

	if len(input.Password) < 8 {
		return nil, ErrPasswordTooShort
	}

	_, err := s.userRepository.GetUserByEmail(ctx, input.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	user := models.User{
		Email:        input.Email,
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	createdUser, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &createdUser, nil
}

func (u *UserService) Login(ctx context.Context, loginInput models.LoginInput) (string, error) {
	if loginInput.Email == "" {
		return "", ErrInvalidCredentials
	}

	if loginInput.Password == "" {
		return "", ErrInvalidCredentials
	}

	user, err := u.userRepository.GetUserByEmail(ctx, loginInput.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(loginInput.Password),
	)

	if err != nil {
		return "", ErrInvalidCredentials
	}

	tokenString, err := auth.GenerateToken(user, u.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (u *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return u.userRepository.GetUsers(ctx)
}

func (u *UserService) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	return u.userRepository.GetUserByID(ctx, userID)
}
