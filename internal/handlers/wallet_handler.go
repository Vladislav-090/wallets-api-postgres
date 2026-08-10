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


// @Summary Create Wallet
// @Description Creates a wallet for the authenticated user
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.WalletInput true "Wallet data"
// @Success 201 {object} models.Wallet
// @Failure 400 {object} response.ResponseError
// @Failure 401 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /wallets [post]
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

	createdWallet, err := h.walletService.CreateWallet(r.Context(), userID, input)
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


// @Summary Get wallets
// @Description Returns wallets of the authenticated user
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Wallet
// @Failure 401 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /wallets [get]
func (h *WalletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "failed to get userID")
		return
	}

	wallets, err := h.walletService.GetWallets(r.Context(), claims.UserID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get wallets")
		return
	}

	response.WriteJSON(w, http.StatusOK, wallets)
}


// @Summary Get wallet by ID
// @Description Returns a wallet of the authenticated user by id
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wallet ID"
// @Success 200 {object} models.Wallet
// @Failure 401 {object} response.ResponseError
// @Failure 400 {object} response.ResponseError
// @Failure 404 {object} response.ResponseError
// @Router /wallets/{id} [get]
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

	wallet, err := h.walletService.GetWalletByID(r.Context(), idInt, claims.UserID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "wallet not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, wallet)
}

// @Summary Update wallet
// @Description Update wallet name for the authenticated user
// @Tags wallets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wallet ID"
// @Param input body models.UpdateWalletInput true "Wallet update data"
// @Success 200 {object} models.Wallet
// @Failure 401 {object} response.ResponseError
// @Failure 400 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /wallets/{id} [patch]
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

	updatedWallet, err := h.walletService.UpdateWallet(r.Context(), idInt, claims.UserID, input.Name)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to update wallet")
		return
	}
	response.WriteJSON(w, http.StatusOK, updatedWallet)
}

// @Summary Delete wallet
// @Description Deletes a wallet of the authenticated user
// @Tags wallets
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wallet ID"
// @Success 200 {object} response.ResponseSuccess
// @Failure 400 {object} response.ResponseError
// @Failure 401 {object} response.ResponseError
// @Failure 404 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /wallets/{id} [delete]
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

	err = h.walletService.DeleteWallet(r.Context(), idInt, claims.UserID)
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

// @Summary Get all wallets
// @Description Returns all wallets. Admin only.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Wallet
// @Failure 401 {object} response.ResponseError
// @Failure 403 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /admin/wallets [get]
func (h *WalletHandler) GetAllWallets(w http.ResponseWriter, r *http.Request) {
	wallets, err := h.walletService.GetAllWallets(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get wallets")
		return
	}

	response.WriteJSON(w, http.StatusOK, wallets)
}


// @Summary Get wallets by user ID
// @Description Returns all wallets for a specific user. Admin only.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {array} models.Wallet
// @Failure 400 {object} response.ResponseError
// @Failure 401 {object} response.ResponseError
// @Failure 403 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /admin/users/{id}/wallets [get]
func (h *WalletHandler) GetWalletsByUserID(w http.ResponseWriter, r *http.Request) {

	idParam := r.PathValue("id")

	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || userID <= 0 {
		response.WriteError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	wallets, err := h.walletService.GetWalletsByUserID(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get wallets")
		return
	}

	response.WriteJSON(w, http.StatusOK, wallets)
}


// @Summary Get wallet by ID
// @Description Returns any wallet by ID. Admin only.
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wallet ID"
// @Success 200 {object} models.Wallet
// @Failure 400 {object} response.ResponseError
// @Failure 401 {object} response.ResponseError
// @Failure 403 {object} response.ResponseError
// @Failure 404 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /admin/wallets/{id} [get]
func (h *WalletHandler) GetWalletByIDForAdmin(w http.ResponseWriter, r *http.Request) {

	idParam := r.PathValue("id")

	walletID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || walletID <= 0 {
		response.WriteError(w, http.StatusBadRequest, "invalid wallet ID")
		return
	}

	wallet, err := h.walletService.GetWalletByIDForAdmin(r.Context(), walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "wallet not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to get wallet")
		return
	}

	response.WriteJSON(w, http.StatusOK, wallet)
}
