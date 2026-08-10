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

// @Summary Create user 
// @Description Registers a new user
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.RegisterInput true "User registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} response.ResponseError
// @Failure 409 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /register [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input models.RegisterInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	createdUser, err := h.userService.CreateUser(r.Context(), input)
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


// @Summary Login user
// @Description Authenticates user and returns JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param input body models.LoginInput true "User Login data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} response.ResponseError
// @Failure 401 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

	var input models.LoginInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	tokenString, err := h.userService.Login(r.Context(), input)
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

// @Summary Get users
// @Description Returns all users
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.User
// @Failure 401 {object} response.ResponseError
// @Failure 403 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router  /admin/users [get]
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetUsers(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "failed to get users")
		return
	}

	response.WriteJSON(w, http.StatusOK, users)
}

// @Summary Get user by id
// @Description Returns a user by ID. Admin only
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} response.ResponseError
// @Failure 401 {object} response.ResponseError
// @Failure 403 {object} response.ResponseError
// @Failure 404 {object} response.ResponseError
// @Failure 500 {object} response.ResponseError
// @Router  /admin/users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idParam := r.PathValue("id")
	userID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || userID <= 0 {
		response.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), userID)
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
