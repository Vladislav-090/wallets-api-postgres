package handlers

import (
	"net/http"
	"wallets-api-postgres/internal/response"
)

// @Summary Health check
// @Description Checks that API is running
// @Tags health
// @Produce json
// @Success 200 {object} response.ResponseSuccess
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	result := response.ResponseSuccess{
		Message: "wallets-api-postgres is running",
	}
	response.WriteJSON(w, http.StatusOK, result)
}
