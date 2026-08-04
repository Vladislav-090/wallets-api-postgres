package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"wallets-api-postgres/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	getUserByEmailResult models.User
	getUserByEmailError  error
	getUserByEmailCalled bool

	createUserCalled bool
	receivedUser     models.User
	createUserResult models.User
	createUserError  error

	getUserByIDResult models.User
	getUserByIDError  error
	getUserByIDCalled bool

	getUsersResult []models.User
	getUsersError  error
	getUsersCalled bool
}

func (f *fakeUserRepository) GetUserByEmail(ctx context.Context, userEmail string) (models.User, error) {
	f.getUserByEmailCalled = true
	return f.getUserByEmailResult, f.getUserByEmailError
}

func (f *fakeUserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	f.createUserCalled = true
	f.receivedUser = user

	return f.createUserResult, f.createUserError
}

func (f *fakeUserRepository) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	f.getUserByIDCalled = true
	return f.getUserByIDResult, f.getUserByIDError
}

func (f *fakeUserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	f.getUsersCalled = true
	return f.getUsersResult, f.getUsersError
}

func TestUserService_CreateUser_Validation(t *testing.T) {
	tests := []struct {
		name          string
		input         models.RegisterInput
		expectedError error
	}{
		{name: "email is required",
			input: models.RegisterInput{
				Email:    "",
				Password: "password123",
			},
			expectedError: ErrEmailRequired,
		},

		{name: "password is required",
			input: models.RegisterInput{
				Email:    "test@gmail.com",
				Password: "",
			},
			expectedError: ErrPasswordRequired,
		},
		{name: "password too short",
			input: models.RegisterInput{
				Email:    "test@gmail.com",
				Password: "123",
			},
			expectedError: ErrPasswordTooShort,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeRepo := &fakeUserRepository{}
			userService := NewUserService(fakeRepo, "test-secret")
			_, err := userService.CreateUser(context.Background(), test.input)
			if !errors.Is(err, test.expectedError) {
				t.Fatalf("expected %v, got %v", test.expectedError, err)
			}

			if fakeRepo.createUserCalled {
				t.Fatal("repository CreateUser must not be called")
			}
		})
	}

}

func TestUserService_CreateUser_RepositoryError(t *testing.T) {
	repErr := errors.New("repository error")
	fakeRepo := &fakeUserRepository{
		getUserByEmailError: sql.ErrNoRows,
		createUserError:     repErr,
	}

	service := NewUserService(fakeRepo, "test-secret")
	input := models.RegisterInput{
		Email:    "test12345@gmail.com",
		Password: "password123",
	}
	_, err := service.CreateUser(context.Background(), input)

	if !errors.Is(err, repErr) {
		t.Fatalf("expected %v, got %v", repErr, err)
	}

	if !fakeRepo.createUserCalled {
		t.Fatal("repository CreateUser must be called")
	}
}

func TestUserService_CreateUser_EmailAlreadyExists(t *testing.T) {
	fakeRepo := &fakeUserRepository{
		getUserByEmailResult: models.User{
			ID:    1,
			Email: "test@gmail.com",
		},
	}
	service := NewUserService(fakeRepo, "test-secret")
	input := models.RegisterInput{
		Email:    "test@gmail.com",
		Password: "password123",
	}

	_, err := service.CreateUser(context.Background(), input)
	if !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected %v, got %v", ErrEmailAlreadyExists, err)
	}

	if fakeRepo.createUserCalled {
		t.Fatal("repository CreateUser must not be called")
	}
}

func TestUserService_CreateUser_GetUserByEmailError(t *testing.T) {
	repErr := errors.New("get user by email error")

	fakeRepo := &fakeUserRepository{
		getUserByEmailError: repErr,
	}

	service := NewUserService(fakeRepo, "test-secret")

	input := models.RegisterInput{
		Email:    "test@gmail.com",
		Password: "password123",
	}

	_, err := service.CreateUser(context.Background(), input)
	if !errors.Is(err, repErr) {
		t.Fatalf("expected %v, got %v", repErr, err)
	}
	if fakeRepo.createUserCalled {
		t.Fatal("repository CreateUser must not be called")
	}
}

func TestUserService_CreateUser_Success(t *testing.T) {
	fakeRepo := &fakeUserRepository{
		getUserByEmailError: sql.ErrNoRows,
		createUserResult: models.User{
			ID:    1,
			Email: "test@gmail.com",
			Role:  "user",
		},
	}
	service := NewUserService(fakeRepo, "test-secret")

	input := models.RegisterInput{
		Email:    "test@gmail.com",
		Password: "password123",
	}

	createdUser, err := service.CreateUser(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fakeRepo.createUserCalled {
		t.Fatal("repository CreateUser must be called")
	}

	if fakeRepo.receivedUser.Email != input.Email {
		t.Fatalf("expected %s, got %s", input.Email, fakeRepo.receivedUser.Email)
	}

	if fakeRepo.receivedUser.Role != "user" {
		t.Fatalf("expected role user, got %s", fakeRepo.receivedUser.Role)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(fakeRepo.receivedUser.PasswordHash),
		[]byte(input.Password),
	)
	if err != nil {
		t.Fatalf("password was not hashed correctly %v", err)
	}

	if createdUser.ID != fakeRepo.createUserResult.ID {
		t.Fatalf("expected user ID %d, got %d", fakeRepo.createUserResult.ID, createdUser.ID)
	}
}

func TestUserService_Login_Validation(t *testing.T) {
	tests := []struct {
		name          string
		input         models.LoginInput
		expectedError error
	}{
		{
			name: "email is required",
			input: models.LoginInput{
				Email:    "",
				Password: "password123",
			},
			expectedError: ErrInvalidCredentials,
		},
		{
			name: "password is required",
			input: models.LoginInput{
				Email:    "test@gmail.com",
				Password: "",
			},
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeRepo := &fakeUserRepository{}
			service := NewUserService(fakeRepo, "test-secret")
			_, err := service.Login(context.Background(), test.input)
			if !errors.Is(err, test.expectedError) {
				t.Fatalf("expected %v, got %v", test.expectedError, err)
			}

			if fakeRepo.getUserByEmailCalled {
				t.Fatal("repository GetUserByEmail must not be called")
			}
		})
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	fakeRepo := &fakeUserRepository{
		getUserByEmailError: sql.ErrNoRows,
	}
	service := NewUserService(fakeRepo, "test-secret")

	input := models.LoginInput{
		Email:    "test@gmail.com",
		Password: "password123",
	}

	_, err := service.Login(context.Background(), input)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected %v, got %v", ErrInvalidCredentials, err)
	}

	if !fakeRepo.getUserByEmailCalled {
		t.Fatal("repository GetUserByEmail must be called")
	}
}

func TestUserService_Login_GetUserByEmailError(t *testing.T) {
	repoErr := errors.New("get user by email error")
	fakeRepo := &fakeUserRepository{
		getUserByEmailError: repoErr,
	}
	service := NewUserService(fakeRepo, "test-secret")

	input := models.LoginInput{
		Email:    "test@gmail.com",
		Password: "password123",
	}

	_, err := service.Login(context.Background(), input)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}

	if !fakeRepo.getUserByEmailCalled {
		t.Fatal("repository GetUserByEmail must be called")
	}
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("correct-password"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	fakeRepo := &fakeUserRepository{
		getUserByEmailResult: models.User{
			ID:           1,
			Email:        "test@gmail.com",
			PasswordHash: string(passwordHash),
			Role:         "user",
		},
	}

	service := NewUserService(fakeRepo, "test-secret")

	input := models.LoginInput{
		Email:    "test@gmail.com",
		Password: "wrong-password",
	}

	token, err := service.Login(context.Background(), input)
	if token != "" {
		t.Fatalf("expected empty token, got %s", token)
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected %v, got %v", ErrInvalidCredentials, err)
	}
	if !fakeRepo.getUserByEmailCalled {
		t.Fatal("repository GetUserByEmail must be called")
	}
}

func TestUserService_Login_Success(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("failed to generate password hash: %v", err)
	}

	fakeRepo := &fakeUserRepository{
		getUserByEmailResult: models.User{
			ID:           1,
			Email:        "test@gmail.com",
			PasswordHash: string(passwordHash),
			Role:         "user",
		},
	}

	service := NewUserService(fakeRepo, "test-secret")

	input := models.LoginInput{
		Email:    "test@gmail.com",
		Password: "password123",
	}

	token, err := service.Login(context.Background(), input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if token == "" {
		t.Fatalf("expected not empty token, got %s", token)
	}

	if !fakeRepo.getUserByEmailCalled {
		t.Fatal("repository GetUserByEmail must be called")
	}
}

func TestUserService_GetUserByID_RepositoryError(t *testing.T) {
	repoErr := errors.New("failed to get user")
	fakeRepo := &fakeUserRepository{
		getUserByIDError: repoErr,
	}

	service := NewUserService(fakeRepo, "test-secret")

	userID := int64(1)

	_, err := service.GetUserByID(context.Background(), userID)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected %v, got %v", repoErr, err)
	}

	if !fakeRepo.getUserByIDCalled {
		t.Fatalf("repository GetUserByID must be called")
	}
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	expectedUser := models.User{
		ID:    1,
		Email: "test@gmail.com",
		Role:  "user",
	}

	fakeRepo := &fakeUserRepository{
		getUserByIDResult: expectedUser,
	}

	service := NewUserService(fakeRepo, "test-secret")

	actualUser, err := service.GetUserByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actualUser.ID != expectedUser.ID {
		t.Fatalf("expected user ID %d, got %d", expectedUser.ID, actualUser.ID)
	}

	if actualUser.Email != expectedUser.Email {
		t.Fatalf("expected email %s, got %s", expectedUser.Email, actualUser.Email)
	}

	if actualUser.Role != expectedUser.Role {
		t.Fatalf("expected role %s, got %s", expectedUser.Role, actualUser.Role)
	}

	if !fakeRepo.getUserByIDCalled {
		t.Fatal("repository GetUserByID must be called")
	}
}

func TestUserService_GetUsers_RepositoryError(t *testing.T) {
	repErr := errors.New("failed to get users")
	fakeRepo := &fakeUserRepository{
		getUsersError: repErr,
	}

	service := NewUserService(fakeRepo, "test-secret")

	_, err := service.GetUsers(context.Background())
	if !errors.Is(err, repErr) {
		t.Fatalf("expected %v, got %v", repErr, err)
	}

	if !fakeRepo.getUsersCalled {
		t.Fatal("repository GetUsers must be called")
	}
}

func TestUserService_GetUsers_Success(t *testing.T) {
	expectedUsers := []models.User{
		{
			ID:    1,
			Email: "test1@gmail.com",
			Role:  "user",
		},
		{
			ID:    2,
			Email: "test2@gmail.com",
			Role:  "user",
		},
	}

	fakeRepo := &fakeUserRepository{
		getUsersResult: expectedUsers,
	}

	service := NewUserService(fakeRepo, "test-secret")

	actualUsers, err := service.GetUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(actualUsers) != len(expectedUsers) {
		t.Fatalf("expected %d users, got %d", len(expectedUsers), len(actualUsers))
	}

	for i := range actualUsers {
		if actualUsers[i].Email != expectedUsers[i].Email {
			t.Fatalf(
				"expected email %s, got %s",
				expectedUsers[i].Email,
				actualUsers[i].Email,
			)
		}
		if actualUsers[i].ID != expectedUsers[i].ID {
			t.Fatalf(
				"expected user ID %d, got %d",
				expectedUsers[i].ID,
				actualUsers[i].ID,
			)
		}
	}

	if !fakeRepo.getUsersCalled {
		t.Fatal("repository GetUsers must be called")
	}
}
