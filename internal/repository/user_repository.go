package repository

import (
	"database/sql"
	"wallets-api-postgres/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u *UserRepository) CreateUser(user models.User) (models.User, error) {
	query := `
	INSERT INTO users (email, password_hash, role)
	VALUES ($1, $2, $3)
	RETURNING id, email, password_hash, role, created_at, updated_at`

	err := u.db.QueryRow(query, user.Email, user.PasswordHash, user.Role).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	query := `
	SELECT id, email, password_hash, role, created_at, updated_at
	FROM users
	WHERE email = $1
`

	err := u.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
