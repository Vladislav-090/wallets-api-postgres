package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"wallets-api-postgres/internal/middleware"
	"wallets-api-postgres/internal/models"
	"wallets-api-postgres/internal/repository"
	"wallets-api-postgres/internal/response"
	"wallets-api-postgres/internal/service"

	"github.com/shopspring/decimal"
)

type TransferHandler struct {
	transferService *service.TransferService
}

func NewTransferHandler(transferService *service.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

func (h *TransferHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "failed to get userID")
		return
	}

	var input models.TransferInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if input.FromWalletID == 0 {
		response.WriteError(w, http.StatusBadRequest, "fromWallet is empty")
		return
	}

	if input.ToWalletID == 0 {
		response.WriteError(w, http.StatusBadRequest, "toWallet is empty")
		return
	}

	if input.Amount.LessThanOrEqual(decimal.Zero) {
		response.WriteError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}

	transfer, err := h.transferService.CreateTransfer(
		claims.UserID,
		input.FromWalletID,
		input.ToWalletID,
		input.Amount,
	)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrSameWallets):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, repository.ErrTooSmallBalance):
			response.WriteError(w, http.StatusConflict, err.Error())
		case errors.Is(err, repository.ErrWalletsCurrencies):
			response.WriteError(w, http.StatusConflict, err.Error())
		case errors.Is(err, repository.ErrZeroAmount):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, sql.ErrNoRows):
			response.WriteError(w, http.StatusNotFound, "wallet not found")

		default:
			response.WriteError(w, http.StatusInternalServerError, "failed to create transfer")

		}
		return
	}
	response.WriteJSON(w, http.StatusCreated, transfer)
}
