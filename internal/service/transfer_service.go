package service

import (
	"wallets-api-postgres/internal/models"
	"wallets-api-postgres/internal/repository"

	"github.com/shopspring/decimal"
)

type TransferService struct {
	transferRepository *repository.TransferRepository
}

func NewTransferService(transferRepository *repository.TransferRepository) *TransferService {
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
