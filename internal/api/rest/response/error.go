package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: message})
}

func BadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, message)
}

func Unauthorized(w http.ResponseWriter) {
	WriteError(w, http.StatusUnauthorized, "unauthorized")
}

func NotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, message)
}

func Conflict(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusConflict, message)
}

func Forbidden(w http.ResponseWriter) {
	WriteError(w, http.StatusForbidden, "forbidden")
}

func InternalError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "internal server error")
}
