package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var input models.LoginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	tokenString, err := h.userService.Login(input)
	if errors.Is(err, service.ErrInvalidCredentials) {
		response.WriteError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
	response.WriteJSON(w, http.StatusOK, map[string]string{
		"token": tokenString,
	})

}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetUsers()
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get users")
		return
	}

	response.WriteJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || userID <= 0 {
		response.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			response.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "failed to get user")
		return

	}
	response.WriteJSON(w, http.StatusOK, user)
}
