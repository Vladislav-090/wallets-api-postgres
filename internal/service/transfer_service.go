package service

import (
	"context"
	"errors"
	"wallets-api-postgres/internal/models"

	"github.com/shopspring/decimal"
)

type TransferRepository interface {
	CreateTransfer(ctx context.Context, userID int64, fromWalletID int64, toWalletID int64, amount decimal.Decimal) (models.Transfer, error)
	GetTransfers(userID int64) ([]models.Transfer, error)
	GetTransferByID(transferID int64, userID int64) (models.Transfer, error)
	GetAllTransfers() ([]models.Transfer, error)
	GetTransferByIDForAdmin(transferID int64) (models.Transfer, error)
}

type TransferService struct {
	transferRepository TransferRepository
}

func NewTransferService(transferRepository TransferRepository) *TransferService {
	return &TransferService{
		transferRepository: transferRepository,
	}
}

var (
	ErrInvalidTransferUserID = errors.New("invalid user id")
	ErrInvalidFromWalletID   = errors.New("invalid from wallet id")
	ErrInvalidToWalletID     = errors.New("invalid to wallet id")
	ErrInvalidTransferAmount = errors.New("amount must be greater than zero")
	ErrSameWallet            = errors.New("cannot transfer to the same wallet")
)

func (s *TransferService) CreateTransfer(
	ctx context.Context,
	userID int64,
	fromWalletID int64,
	toWalletID int64,
	amount decimal.Decimal,
) (models.Transfer, error) {

	if userID <= 0 {
		return models.Transfer{}, ErrInvalidTransferUserID
	}
	if fromWalletID <= 0 {
		return models.Transfer{}, ErrInvalidFromWalletID
	}

	if toWalletID <= 0 {
		return models.Transfer{}, ErrInvalidToWalletID
	}

	if fromWalletID == toWalletID {
		return models.Transfer{}, ErrSameWallet
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return models.Transfer{}, ErrInvalidTransferAmount
	}

	return s.transferRepository.CreateTransfer(
		ctx,
		userID,
		fromWalletID,
		toWalletID,
		amount,
	)
}

func (s *TransferService) GetTransfers(userID int64) ([]models.Transfer, error) {
	return s.transferRepository.GetTransfers(userID)
}

func (s *TransferService) GetTransferByID(transferID int64, userID int64) (models.Transfer, error) {
	return s.transferRepository.GetTransferByID(transferID, userID)
}

func (s *TransferService) GetAllTransfers() ([]models.Transfer, error) {
	return s.transferRepository.GetAllTransfers()
}

func (s *TransferService) GetTransferByIDForAdmin(transferID int64) (models.Transfer, error) {
	return s.transferRepository.GetTransferByIDForAdmin(transferID)
}
