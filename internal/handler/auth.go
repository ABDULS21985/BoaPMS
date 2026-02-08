package handler

import (
	"encoding/json"
	"net/http"

	"github.com/enterprise-pms/pms-api/internal/domain/auth"
	"github.com/enterprise-pms/pms-api/internal/middleware"
	"github.com/enterprise-pms/pms-api/internal/service"
	"github.com/enterprise-pms/pms-api/pkg/response"
	"github.com/rs/zerolog"
)

// AuthHandler handles authentication HTTP endpoints.
// Mirrors the .NET AuthController with login and refresh endpoints.
type AuthHandler struct {
	authSvc service.AuthService
	log     zerolog.Logger
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(authSvc service.AuthService, log zerolog.Logger) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, log: log}
}

// Login handles POST /api/v1/auth/login
// Mirrors .NET AuthController.Login — dual-mode AD/local authentication.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.AuthenticateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	result, err := h.authSvc.AuthenticateAD(r.Context(), req.Username, req.Password)
	if err != nil {
		h.log.Warn().Err(err).Str("user", req.Username).Msg("Authentication failed")
		response.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.OK(w, result)
}

// RefreshToken handles POST /api/v1/auth/refresh
// Exchanges a refresh token for a new access token.
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req auth.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		response.Error(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	newAccessToken, err := h.authSvc.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	response.OK(w, map[string]string{
		"access_token": newAccessToken,
	})
}

// ValidateToken handles GET /api/v1/auth/validate
// Validates the current bearer token and returns the claims.
// The JWT middleware already validated the token before this handler runs.
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	userID, _ := r.Context().Value(middleware.UserIDKey).(string)
	email, _ := r.Context().Value(middleware.EmailKey).(string)
	roles, _ := r.Context().Value(middleware.RolesKey).([]string)
	perms, _ := r.Context().Value(middleware.PermissionsKey).([]string)
	ou, _ := r.Context().Value(middleware.OrganizationalUnitKey).(string)

	claims := auth.CurrentUserData{
		UserID:             userID,
		Email:              email,
		Roles:              roles,
		Permissions:        perms,
		OrganizationalUnit: ou,
	}
	response.OK(w, claims)
}
