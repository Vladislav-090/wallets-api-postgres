package repository

import (
	"database/sql"
	"wallets-api-postgres/internal/models"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (w *WalletRepository) CreateWallet(wallet models.Wallet) (models.Wallet, error) {
	query := `INSERT INTO wallets (user_id, name, currency, balance)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, name, currency, balance, created_at, updated_at`

	err := w.db.QueryRow(query, wallet.UserID, wallet.Name, wallet.Currency, wallet.Balance).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Name,
		&wallet.Currency,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (w *WalletRepository) GetWallets(userID int64) ([]models.Wallet, error) {
	query := `SELECT id, user_id, name, currency, balance, created_at, updated_at
	FROM wallets
	WHERE user_id = $1
	ORDER BY created_at DESC`

	rows, err := w.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	wallets := make([]models.Wallet, 0)

	for rows.Next() {
		var wallet models.Wallet

		err := rows.Scan(
			&wallet.ID,
			&wallet.UserID,
			&wallet.Name,
			&wallet.Currency,
			&wallet.Balance,
			&wallet.CreatedAt,
			&wallet.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (w *WalletRepository) GetWalletByID(walletID int64, userID int64) (models.Wallet, error) {
	query := `SELECT id, user_id, name, currency, balance, created_at, updated_at
	FROM wallets
	WHERE id = $1 AND user_id = $2`

	var wallet models.Wallet

	err := w.db.QueryRow(query,
		walletID,
		userID,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Name,
		&wallet.Currency,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (w *WalletRepository) UpdateWallet(walletID int64, userID int64, name string) (models.Wallet, error) {
	query := `
	UPDATE wallets
	SET name = $1,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $2 AND user_id = $3
	RETURNING id, user_id, name, currency, balance, created_at, updated_at`

	var wallet models.Wallet

	err := w.db.QueryRow(query,
		name,
		walletID,
		userID,
	).Scan(
		&wallet.ID,
		&wallet.UserID,
		&wallet.Name,
		&wallet.Currency,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (w *WalletRepository) DeleteWallet(walletID int64, userID int64) error {
	query := `
	DELETE FROM wallets
	WHERE id = $1 AND user_id = $2
	`
	result, err := w.db.Exec(query, walletID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
