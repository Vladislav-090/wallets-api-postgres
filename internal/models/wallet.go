package models

import "time"

type Wallet struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WalletInput struct {
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

type UpdateWalletInput struct {
	Name string `json:"name"`
}
