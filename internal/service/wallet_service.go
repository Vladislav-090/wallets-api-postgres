package service

import (
	"errors"
	"wallets-api-postgres/internal/models"

	"github.com/shopspring/decimal"
)

type WalletRepository interface {
	CreateWallet(wallet models.Wallet) (models.Wallet, error)
	GetWallets(userID int64) ([]models.Wallet, error)
	GetWalletByID(walletID int64, userID int64) (models.Wallet, error)
	UpdateWallet(walletID int64, userID int64, name string) (models.Wallet, error)
	DeleteWallet(walletID int64, userID int64) error
	GetAllWallets() ([]models.Wallet, error)
	GetWalletsByUserID(userID int64) ([]models.Wallet, error)
	GetWalletByIDForAdmin(walletID int64) (models.Wallet, error)
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

func (w *WalletService) CreateWallet(userID int64, input models.WalletInput) (*models.Wallet, error) {
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

	createdWallet, err := w.repo.CreateWallet(wallet)
	if err != nil {
		return nil, err
	}

	return &createdWallet, nil

}

func (w *WalletService) GetWallets(userID int64) ([]models.Wallet, error) {

	wallets, err := w.repo.GetWallets(userID)
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (w *WalletService) GetWalletByID(walletID int64, userID int64) (models.Wallet, error) {
	wallet, err := w.repo.GetWalletByID(walletID, userID)
	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}

func (w *WalletService) UpdateWallet(walletID int64, userID int64, name string) (models.Wallet, error) {
	updatedWallet, err := w.repo.UpdateWallet(walletID, userID, name)
	if err != nil {
		return models.Wallet{}, err
	}

	return updatedWallet, nil
}

func (w *WalletService) DeleteWallet(walletID int64, userID int64) error {
	err := w.repo.DeleteWallet(walletID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (w *WalletService) GetAllWallets() ([]models.Wallet, error) {
	wallets, err := w.repo.GetAllWallets()
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (w *WalletService) GetWalletsByUserID(userID int64) ([]models.Wallet, error) {
	wallets, err := w.repo.GetWalletsByUserID(userID)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

func (w *WalletService) GetWalletByIDForAdmin(walletID int64) (models.Wallet, error) {
	wallet, err := w.repo.GetWalletByIDForAdmin(walletID)
	if err != nil {
		return models.Wallet{}, err
	}

	return wallet, nil
}
