package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseSuccess struct {
	Message string `json:"message"`
}

type ResponseError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("failed to encode response", err)
	}
}

func WriteError(w http.ResponseWriter, status int, message string) {
	errorResponse := ResponseError{
		Error: message,
	}
	WriteJSON(w, status, errorResponse)
}
