package repository

import (
	"database/sql"
	"errors"
	"wallets-api-postgres/internal/models"

	"github.com/shopspring/decimal"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{
		db: db,
	}
}

var (
	ErrSameWallets       = errors.New("cannot transfer to the same wallet")
	ErrZeroAmount        = errors.New("amount must be greater than zero")
	ErrWalletsCurrencies = errors.New("wallets has different currencies")
	ErrTooSmallBalance   = errors.New("dont have enough money")
)

func (t *TransferRepository) CreateTransfer(userID int64,
	fromWalletID int64,
	toWalletID int64,
	amount decimal.Decimal,
) (models.Transfer, error) {
	tx, err := t.db.Begin()
	if err != nil {
		return models.Transfer{}, err
	}
	defer tx.Rollback()

	var fromBalance decimal.Decimal
	var fromCurrency string

	query := `
	SELECT balance, currency
	FROM wallets
	WHERE id = $1 AND user_id = $2
	FOR UPDATE`

	err = tx.QueryRow(query, fromWalletID, userID).Scan(
		&fromBalance,
		&fromCurrency,
	)
	if err != nil {
		return models.Transfer{}, err
	}

	var toCurrency string

	toQuery := `
	SELECT currency
	FROM wallets
	WHERE id = $1
	FOR UPDATE`

	err = tx.QueryRow(toQuery, toWalletID).Scan(
		&toCurrency,
	)
	if err != nil {
		return models.Transfer{}, err
	}

	if fromWalletID == toWalletID {
		return models.Transfer{}, ErrSameWallets
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return models.Transfer{}, ErrZeroAmount
	}

	if fromCurrency != toCurrency {
		return models.Transfer{}, ErrWalletsCurrencies
	}

	if fromBalance.LessThan(amount) {
		return models.Transfer{}, ErrTooSmallBalance
	}

	debitQuery := `
	UPDATE wallets
	SET balance = balance - $1,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $2`

	_, err = tx.Exec(debitQuery, amount, fromWalletID)
	if err != nil {
		return models.Transfer{}, err
	}

	creditQuery := `
	UPDATE wallets
	SET balance = balance + $1,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $2`

	_, err = tx.Exec(creditQuery, amount, toWalletID)
	if err != nil {
		return models.Transfer{}, err
	}

	var transfer models.Transfer

	transferQuery := `
	INSERT INTO transfers(
	 from_wallet_id,
	 to_wallet_id,
	 amount,
	 status
	 )
	 VALUES ($1, $2, $3, $4)
	 RETURNING 
	 id,
	 from_wallet_id,
	 to_wallet_id,
	 amount,
	 status,
	 created_at`

	err = tx.QueryRow(
		transferQuery,
		fromWalletID,
		toWalletID,
		amount,
		models.TransferStatusCompleted,
	).Scan(
		&transfer.ID,
		&transfer.FromWalletID,
		&transfer.ToWalletID,
		&transfer.Amount,
		&transfer.Status,
		&transfer.CreatedAt,
	)
	if err != nil {
		return models.Transfer{}, err
	}

	err = tx.Commit()
	if err != nil {
		return models.Transfer{}, err
	}

	return transfer, nil
}
