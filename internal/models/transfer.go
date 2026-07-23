package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransferStatus string

const (
	TransferStatusPending   TransferStatus = "pending"
	TransferStatusCompleted TransferStatus = "completed"
	TransferStatusFailed    TransferStatus = "failed"
)

type Transfer struct {
	ID           int64           `json:"id"`
	FromWalletID int64           `json:"from_wallet_id"`
	ToWalletID   int64           `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
	Status       TransferStatus  `json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
}

type TransferInput struct {
	FromWalletID int64           `json:"from_wallet_id"`
	ToWalletID   int64           `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
}
