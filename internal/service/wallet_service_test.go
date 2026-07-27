package service

import (
	"errors"
	"testing"
	"wallets-api-postgres/internal/models"
)

type fakeWalletRepository struct {
	createWalletCalled bool
	receivedWallet     models.Wallet
	createWalletResult models.Wallet
	createWalletError  error
}

func (f *fakeWalletRepository) CreateWallet(wallet models.Wallet) (models.Wallet, error) {
	f.createWalletCalled = true
	f.receivedWallet = wallet

	return f.createWalletResult, f.createWalletError
}

func (f *fakeWalletRepository) GetWallets(userID int64) ([]models.Wallet, error) {
	return nil, nil
}

func (f *fakeWalletRepository) GetWalletByID(walletID int64, userID int64) (models.Wallet, error) {
	return models.Wallet{}, nil
}

func (f *fakeWalletRepository) UpdateWallet(walletID int64, userID int64, name string) (models.Wallet, error) {
	return models.Wallet{}, nil
}

func (f *fakeWalletRepository) DeleteWallet(walletID int64, userID int64) error {
	return nil
}

func (f *fakeWalletRepository) GetAllWallets() ([]models.Wallet, error) {
	return nil, nil
}

func (f *fakeWalletRepository) GetWalletsByUserID(userID int64) ([]models.Wallet, error) {
	return nil, nil
}

func (f *fakeWalletRepository) GetWalletByIDForAdmin(walletID int64) (models.Wallet, error) {
	return models.Wallet{}, nil
}

func TestWalletService_CreateWallet_Success(t *testing.T) {
	fakeRepo := &fakeWalletRepository{
		createWalletResult: models.Wallet{
			ID:       10,
			UserID:   1,
			Name:     "Main",
			Currency: "USD",
		},
	}
	service := NewWalletService(fakeRepo)

	result, err := service.CreateWallet(1, models.WalletInput{
		Name:     "Main",
		Currency: "USD",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !fakeRepo.createWalletCalled {
		t.Fatal("expected repository CreateWallet to be called")
	}
	if result == nil {
		t.Fatal("expected wallet result, got nil")
	}
	if result.ID != 10 {
		t.Fatalf("expected wallet ID 10, got %d", result.ID)
	}

}

func TestWalletService_CreateWallet_RepositoryError(t *testing.T) {
	repErr := errors.New("database error")

	fakeRepo := &fakeWalletRepository{
		createWalletError: repErr,
	}

	service := NewWalletService(fakeRepo)

	_, err := service.CreateWallet(1, models.WalletInput{
		Name:     "Main",
		Currency: "USD",
	})

	if !errors.Is(err, repErr) {
		t.Fatalf("expected repository error, got %v", err)
	}

	if !fakeRepo.createWalletCalled {
		t.Fatal("expected repository CreateWallet to be called")
	}
}

func TestWalletService_CreateWallet_Validation(t *testing.T) {
	tests := []struct {
		name        string
		input       models.WalletInput
		expectedErr error
	}{
		{
			name: "name is required",
			input: models.WalletInput{
				Name:     "",
				Currency: "USD",
			},
			expectedErr: ErrNameRequired,
		},
		{
			name: "currency is required",
			input: models.WalletInput{
				Name:     "Main",
				Currency: "",
			},
			expectedErr: ErrCurrencyRequired,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fakeRepo := &fakeWalletRepository{}
			walletService := NewWalletService(fakeRepo)
			_, err := walletService.CreateWallet(1, test.input)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("expected %v, got %v", test.expectedErr, err)
			}

			if fakeRepo.createWalletCalled {
				t.Fatal("repository CreateWallet must not be called")
			}
		})

	}
}
