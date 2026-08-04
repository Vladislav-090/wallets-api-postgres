package service

import (
	"context"
	"errors"
	"wallets-api-postgres/internal/models"

	"github.com/shopspring/decimal"
)

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet models.Wallet) (models.Wallet, error)
	GetWallets(ctx context.Context, userID int64) ([]models.Wallet, error)
	GetWalletByID(ctx context.Context, walletID int64, userID int64) (models.Wallet, error)
	UpdateWallet(ctx context.Context, walletID int64, userID int64, name string) (models.Wallet, error)
	DeleteWallet(ctx context.Context, walletID int64, userID int64) error
	GetAllWallets(ctx context.Context) ([]models.Wallet, error)
	GetWalletsByUserID(ctx context.Context, userID int64) ([]models.Wallet, error)
	GetWalletByIDForAdmin(ctx context.Context, walletID int64) (models.Wallet, error)
}

type WalletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}

var (
	ErrNameRequired     = errors.New("name is required")
	ErrCurrencyRequired = errors.New("currency is required")
)

func (w *WalletService) CreateWallet(ctx context.Context, userID int64, input models.WalletInput) (*models.Wallet, error) {
	if input.Name == "" {
		return nil, ErrNameRequired
	}

	if input.Currency == "" {
		return nil, ErrCurrencyRequired
	}

	wallet := models.Wallet{
		UserID:   userID,
		Name:     input.Name,
		Currency: input.Currency,
		Balance:  decimal.Zero,
	}

	createdWallet, err := w.repo.CreateWallet(ctx, wallet)
	if err != nil {
		return nil, err
	}

	return &createdWallet, nil

}

func (w *WalletService) GetWallets(ctx context.Context, userID int64) ([]models.Wallet, error) {

	wallets, err := w.repo.GetWallets(ctx, userID)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (w *WalletService) GetWalletByID(ctx context.Context, walletID int64, userID int64) (models.Wallet, error) {
	wallet, err := w.repo.GetWalletByID(ctx, walletID, userID)
	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (w *WalletService) UpdateWallet(ctx context.Context, walletID int64, userID int64, name string) (models.Wallet, error) {
	updatedWallet, err := w.repo.UpdateWallet(ctx, walletID, userID, name)
	if err != nil {
		return models.Wallet{}, err
	}

	return updatedWallet, nil
}

func (w *WalletService) DeleteWallet(ctx context.Context, walletID int64, userID int64) error {
	err := w.repo.DeleteWallet(ctx, walletID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (w *WalletService) GetAllWallets(ctx context.Context) ([]models.Wallet, error) {
	wallets, err := w.repo.GetAllWallets(ctx)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (w *WalletService) GetWalletsByUserID(ctx context.Context, userID int64) ([]models.Wallet, error) {
	wallets, err := w.repo.GetWalletsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (w *WalletService) GetWalletByIDForAdmin(ctx context.Context, walletID int64) (models.Wallet, error) {
	wallet, err := w.repo.GetWalletByIDForAdmin(ctx, walletID)
	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}
