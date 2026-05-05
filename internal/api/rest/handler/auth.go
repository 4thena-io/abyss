package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/4thena-io/abyss/internal/api/rest/middleware"
	"github.com/4thena-io/abyss/internal/api/rest/response"
	"github.com/4thena-io/abyss/internal/auth"
	"github.com/4thena-io/abyss/internal/service"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		response.InternalError(w)
		return
	}
	state := hex.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((15 * time.Minute).Seconds()),
	})

	http.Redirect(w, r, h.service.AuthURL(state), http.StatusFound)
}

func (h *AuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || stateCookie.Value != r.URL.Query().Get("state") {
		response.BadRequest(w, "invalid state")
		return
	}

	// Consume the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "oauth_state",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	token, err := h.service.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		response.BadRequest(w, "failed to exchange code")
		return
	}

	forgeUser, err := h.service.GetForgeUser(r.Context(), token)
	if err != nil {
		response.InternalError(w)
		return
	}

	user, err := h.service.GetOrCreateUser(r.Context(), forgeUser)
	if errors.Is(err, service.ErrUnauthorized) {
		http.Redirect(w, r, "/login?error=access_denied", http.StatusFound)
		return
	}
	if err != nil {
		response.InternalError(w)
		return
	}

	jwt, err := h.service.GenerateJWT(user)
	if err != nil {
		response.InternalError(w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((24 * time.Hour).Seconds()),
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.ContextKeyUser).(*auth.Claims)
	if claims == nil {
		response.Unauthorized(w)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		response.InternalError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"is_admin":   user.IsAdmin,
		"avatar_url": user.AvatarURL,
		"forge_type": h.service.ForgeType(),
		"forge_host": h.service.ForgeHost(),
	})
}

func (h *AuthHandler) TokenStatus(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.ContextKeyUser).(*auth.Claims)
	if claims == nil {
		response.Unauthorized(w)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		response.InternalError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"has_token": user.Token != nil})
}

func (h *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.ContextKeyUser).(*auth.Claims)
	if claims == nil {
		response.Unauthorized(w)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		response.InternalError(w)
		return
	}
	token, err := h.service.GeneratePersonalToken(r.Context(), user)
	if err != nil {
		response.InternalError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	claims, _ := r.Context().Value(middleware.ContextKeyUser).(*auth.Claims)
	if claims == nil {
		response.Unauthorized(w)
		return
	}
	user, err := h.service.GetUserByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		response.InternalError(w)
		return
	}
	if err := h.service.RevokePersonalToken(r.Context(), user); err != nil {
		response.InternalError(w)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusNoContent)
}
