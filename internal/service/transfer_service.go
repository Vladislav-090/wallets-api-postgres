package service

import (
	"wallets-api-postgres/internal/models"

	"github.com/shopspring/decimal"
)

type TransferRepository interface {
	CreateTransfer(userID int64, fromWalletID int64, toWalletID int64, amount decimal.Decimal) (models.Transfer, error)
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

func (s *TransferService) CreateTransfer(
	userID int64,
	fromWalletID int64,
	toWalletID int64,
	amount decimal.Decimal,
) (models.Transfer, error) {
	return s.transferRepository.CreateTransfer(
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
