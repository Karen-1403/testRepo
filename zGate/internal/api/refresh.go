package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// RefreshRequest represents refresh token request payload
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshResponse represents refresh token response payload
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
}

// handleRefresh handles POST /api/refresh
func (s *Server) handleRefresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Error("invalid refresh request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	// Validate refresh token
	refreshTokenData, err := s.store.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		utils.Logger.Warn("invalid refresh token", "error", err)
		http.Error(w, "invalid or expired refresh token", http.StatusUnauthorized)
		return
	}

	// Get fresh user data with current roles/permissions
	user, err := s.store.GetUser(refreshTokenData.Username)
	if err != nil {
		utils.Logger.Error("failed to get user", "error", err)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// Get role permissions
	rolePerms, err := s.store.GetPermissionsForRoles(user.Roles)
	if err != nil {
		utils.Logger.Error("failed to get role permissions", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Combine permissions
	permissions := append(rolePerms, user.CustomPermissions...)

	// Create user with permissions for token generation
	userWithPerms := &auth.UserWithPermissions{
		Username:    user.Username,
		Roles:       user.Roles,
		Permissions: permissions,
	}

	// Generate new access token
	accessToken, expiresAt, err := auth.GenerateToken(userWithPerms)
	if err != nil {
		utils.Logger.Error("failed to generate access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Token rotation: generate new refresh token
	userAgent := r.Header.Get("User-Agent")
	ipAddress := getClientIP(r)
	newRefreshToken, err := s.store.CreateRefreshToken(
		user.Username,
		userAgent,
		ipAddress,
		auth.RefreshTokenDuration,
	)
	if err != nil {
		utils.Logger.Error("failed to generate new refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Revoke old refresh token (rotation)
	if err := s.store.RevokeRefreshToken(req.RefreshToken); err != nil {
		utils.Logger.Warn("failed to revoke old refresh token", "error", err)
		// Continue anyway - not critical
	}

	utils.Logger.Info("token refreshed", "username", user.Username, "ip", ipAddress)

	// Calculate expires_in
	expiresIn := int(time.Until(expiresAt).Seconds())

	// Return new tokens
	resp := RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
