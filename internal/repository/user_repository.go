package repository

import (
	"context"
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

func (u *UserRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `
	INSERT INTO users (email, password_hash, role)
	VALUES ($1, $2, $3)
	RETURNING id, email, password_hash, role, created_at, updated_at`

	err := u.db.QueryRowContext(ctx, query, user.Email, user.PasswordHash, user.Role).Scan(
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

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	query := `
	SELECT id, email, password_hash, role, created_at, updated_at
	FROM users
	WHERE email = $1
`

	err := u.db.QueryRowContext(ctx, query, email).Scan(
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

func (u *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `
	SELECT id, email, role, created_at, updated_at
	FROM users
	ORDER BY created_at DESC`

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]models.User, 0)

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserRepository) GetUserByID(ctx context.Context, userID int64) (models.User, error) {
	query := `
	SELECT id, email, role, created_at, updated_at
	FROM users
	WHERE id = $1
	`
	var user models.User

	err := u.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil

}
