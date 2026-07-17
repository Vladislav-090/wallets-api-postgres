package handlers

import (
	"net/http"
	"wallets-api-postgres/internal/response"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "method not allowed!")
		return
	}
	result := response.ResponseSuccess{
		Message: "wallets-api-postgres is running",
	}
	response.WriteJSON(w, http.StatusOK, result)
}
