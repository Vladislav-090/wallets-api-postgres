package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"wallets-api-postgres/internal/middleware"
	"wallets-api-postgres/internal/models"
	"wallets-api-postgres/internal/response"
	"wallets-api-postgres/internal/service"
)

type WalletHandler struct {
	walletService *service.WalletService
}

func NewWalletHandler(service *service.WalletService) *WalletHandler {
	return &WalletHandler{
		walletService: service,
	}
}

func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var input models.WalletInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "failed to get userID")
		return
	}

	userID := claims.UserID

	createdWallet, err := h.walletService.CreateWallet(userID, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNameRequired):
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, service.ErrCurrencyRequired):
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to create wallet")
		return
	}

	response.WriteJSON(w, http.StatusCreated, createdWallet)
}

func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "failed to get userID")
		return
	}

	wallets, err := h.walletService.GetWallets(claims.UserID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get wallets")
		return
	}

	response.WriteJSON(w, http.StatusOK, wallets)
}

func (h *WalletHandler) GetWalletByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "failed to get userID")
		return
	}

	idParam := r.PathValue("id")

	idInt, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid wallet id")
		return
	}

	wallet, err := h.walletService.GetWalletByID(idInt, claims.UserID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "wallet not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, wallet)
}

func (h *WalletHandler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	var input models.UpdateWalletInput

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "failed to get userID")
		return
	}

	idParam := r.PathValue("id")

	idInt, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid ID")
		return
	}

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if input.Name == "" {
		response.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	updatedWallet, err := h.walletService.UpdateWallet(idInt, claims.UserID, input.Name)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to update wallet")
		return
	}
	response.WriteJSON(w, http.StatusOK, updatedWallet)
}

func (h *WalletHandler) DeleteWallet(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "failed to get userID")
		return
	}

	idParam := r.PathValue("id")

	idInt, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.walletService.DeleteWallet(idInt, claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "wallet not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to delete wallet")
		return
	}

	response.WriteJSON(w, http.StatusOK, response.ResponseSuccess{
		Message: "wallet deleted successfully",
	})
}
