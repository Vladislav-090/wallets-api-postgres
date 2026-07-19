package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"wallets-api-postgres/internal/models"
	"wallets-api-postgres/internal/response"
	"wallets-api-postgres/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	createdUser, err := h.userService.CreateUser(input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailRequired):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrPasswordRequired):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrPasswordTooShort):
			response.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, service.ErrEmailAlreadyExists):
			response.WriteError(w, http.StatusConflict, err.Error())
		default:
			response.WriteError(w, http.StatusInternalServerError, "failed to create user")
		}
		return
	}

	response.WriteJSON(w, http.StatusCreated, createdUser)
}
