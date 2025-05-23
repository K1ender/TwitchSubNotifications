package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

func WriteJSON(w http.ResponseWriter, status int, data Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func InternalServerError(w http.ResponseWriter) {
	WriteJSON(w, http.StatusInternalServerError, Response{
		Success: false,
		Message: "Internal Server Error",
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	WriteJSON(w, http.StatusBadRequest, Response{
		Success: false,
		Message: message,
	})
}

func Unauthorized(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, Response{
		Success: false,
		Message: "Unauthorized",
	})
}

func OK(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, Response{
		Success: true,
		Message: "Success",
		Data:    data,
	})
}
