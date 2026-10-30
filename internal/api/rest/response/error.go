package response

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/4thena-io/abyss/internal/database"
	"github.com/4thena-io/abyss/internal/integration/forge"
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

func ServiceUnavailable(w http.ResponseWriter, message string) {
	w.Header().Set("Retry-After", "5")
	WriteError(w, http.StatusServiceUnavailable, message)
}

// FromError preserves unexpected failures as 500, but dependency outages are retryable.
func FromError(w http.ResponseWriter, err error) {
	if errors.Is(err, forge.ErrUnavailable) || database.IsTransient(err) || errors.Is(err, context.DeadlineExceeded) {
		ServiceUnavailable(w, "A required service is temporarily unavailable. Please retry shortly.")
		return
	}
	InternalError(w)
}

func InternalError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "internal server error")
}
