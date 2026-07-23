package handlers

import (
	"net/http"
	"wallets-api-postgres/internal/response"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	result := response.ResponseSuccess{
		Message: "wallets-api-postgres is running",
	}
	response.WriteJSON(w, http.StatusOK, result)
}
