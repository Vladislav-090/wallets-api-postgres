package service

import (
	"context"
	"errors"
	"testing"
	"wallets-api-postgres/internal/models"

	"github.com/shopspring/decimal"
)

type fakeTransferRepository struct {
	createTransferCalled bool
	createTransferError  error
	createTransferResult models.Transfer

	receivedUserID       int64
	receivedFromWalletID int64
	receivedToWalletID   int64
	receivedAmount       decimal.Decimal
}

func (f *fakeTransferRepository) CreateTransfer(
	ctx context.Context,
	userID int64,
	fromWalletID int64,
	toWalletID int64,
	amount decimal.Decimal,
) (models.Transfer, error) {
	f.createTransferCalled = true
	f.receivedUserID = userID
	f.receivedFromWalletID = fromWalletID
	f.receivedToWalletID = toWalletID
	f.receivedAmount = amount

	return f.createTransferResult, f.createTransferError
}

func (f *fakeTransferRepository) GetTransfers(ctx context.Context, userID int64) ([]models.Transfer, error) {
	return nil, nil
}

func (f *fakeTransferRepository) GetTransferByID(ctx context.Context, transferID int64, userID int64) (models.Transfer, error) {
	return models.Transfer{}, nil
}

func (f *fakeTransferRepository) GetAllTransfers(ctx context.Context) ([]models.Transfer, error) {
	return []models.Transfer{}, nil
}

func (f *fakeTransferRepository) GetTransferByIDForAdmin(ctx context.Context, transferID int64) (models.Transfer, error) {
	return models.Transfer{}, nil
}

func TestTransferService_CreateTransfer_Validation(t *testing.T) {
	tests := []struct {
		name         string
		userID       int64
		fromWalletID int64
		toWalletID   int64
		amount       decimal.Decimal
		expectedErr  error
	}{
		{
			name:         "invalid user id",
			userID:       0,
			fromWalletID: 1,
			toWalletID:   2,
			amount:       decimal.RequireFromString("100.00"),
			expectedErr:  ErrInvalidTransferUserID,
		},

		{
			name:         "invalid from wallet id",
			userID:       1,
			fromWalletID: 0,
			toWalletID:   2,
			amount:       decimal.RequireFromString("100.00"),
			expectedErr:  ErrInvalidFromWalletID,
		},

		{
			name:         "invalid to wallet id",
			userID:       1,
			fromWalletID: 2,
			toWalletID:   0,
			amount:       decimal.RequireFromString("100.00"),
			expectedErr:  ErrInvalidToWalletID,
		},

		{
			name:         "invalid transfer amount",
			userID:       1,
			fromWalletID: 2,
			toWalletID:   3,
			amount:       decimal.RequireFromString("0"),
			expectedErr:  ErrInvalidTransferAmount,
		},

		{
			name:         "transfer to same wallet",
			userID:       1,
			fromWalletID: 2,
			toWalletID:   2,
			amount:       decimal.RequireFromString("100.00"),
			expectedErr:  ErrSameWallet,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeRepo := &fakeTransferRepository{}
			service := NewTransferService(fakeRepo)

			_, err := service.CreateTransfer(
				context.Background(),
				test.userID,
				test.fromWalletID,
				test.toWalletID,
				test.amount,
			)

			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("expected %v, got %v", test.expectedErr, err)
			}

			if fakeRepo.createTransferCalled {
				t.Fatal("repository CreateTransfer must not be called")
			}
		})
	}
}

func TestTransferService_CreateTransfer_Success(t *testing.T) {
	expectedTransfer := models.Transfer{
		ID:           1,
		FromWalletID: 2,
		ToWalletID:   3,
		Amount:       decimal.RequireFromString("1000.00"),
	}

	fakeRepo := &fakeTransferRepository{
		createTransferResult: expectedTransfer,
	}

	service := NewTransferService(fakeRepo)

	userID := int64(1)
	fromWalletID := int64(2)
	toWalletID := int64(3)
	amount := decimal.RequireFromString("1000.00")

	actualTransfer, err := service.CreateTransfer(
		context.Background(),
		userID,
		fromWalletID,
		toWalletID,
		amount,
	)
	if err != nil {
		t.Fatal("unexpected err:", err)
	}

	if fakeRepo.receivedUserID != userID {
		t.Fatalf("expected %d, got %d", userID, fakeRepo.receivedUserID)
	}
	if fakeRepo.receivedFromWalletID != fromWalletID {
		t.Fatalf(
			"expected from wallet ID %d, got %d",
			fromWalletID,
			fakeRepo.receivedFromWalletID,
		)
	}

	if fakeRepo.receivedToWalletID != toWalletID {
		t.Fatalf(
			"expected to wallet ID %d, got %d",
			toWalletID,
			fakeRepo.receivedToWalletID,
		)
	}

	if !fakeRepo.receivedAmount.Equal(amount) {
		t.Fatalf(
			"expected amount %s, got %s",
			amount.String(),
			fakeRepo.receivedAmount.String(),
		)
	}

	if actualTransfer.ID != expectedTransfer.ID {
		t.Fatalf(
			"expected transfer ID %d, got %d",
			expectedTransfer.ID,
			actualTransfer.ID,
		)
	}
	if !fakeRepo.createTransferCalled {
		t.Fatal("repository CreateTransfer must be called")
	}

}

func TestTransferService_CreateTransfer_RepositoryError(t *testing.T) {
	repErr := errors.New("failed to create transfer")

	fakeRepo := &fakeTransferRepository{
		createTransferError: repErr,
	}

	service := NewTransferService(fakeRepo)

	_, err := service.CreateTransfer(
		context.Background(),
		1,
		2,
		3,
		decimal.RequireFromString("1000.00"),
	)

	if !errors.Is(err, repErr) {
		t.Fatalf("expected %v, got %v", repErr, err)
	}

	if !fakeRepo.createTransferCalled {
		t.Fatal("repository CreateTransfer must be called")
	}
}
