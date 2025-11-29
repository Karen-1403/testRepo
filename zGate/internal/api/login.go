package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// LoginRequest represents login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents login response payload (OAuth2-style)
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	Username     string `json:"username"`
}

// handleLogin handles POST /api/login
func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Error("invalid login request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	utils.Logger.Info("login attempt", "username", req.Username)

	// Authenticate user
	user, err := s.authenticator.Authenticate(req.Username, req.Password)
	if err != nil {
		utils.Logger.Warn("authentication failed", "username", req.Username, "error", err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT access token
	accessToken, expiresAt, err := auth.GenerateToken(user)
	if err != nil {
		utils.Logger.Error("failed to generate access token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Generate refresh token
	userAgent := r.Header.Get("User-Agent")
	ipAddress := getClientIP(r)
	refreshToken, err := s.store.CreateRefreshToken(
		user.Username,
		userAgent,
		ipAddress,
		auth.RefreshTokenDuration,
	)
	if err != nil {
		utils.Logger.Error("failed to generate refresh token", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("login successful", "username", req.Username, "ip", ipAddress)

	// Calculate expires_in (seconds until expiration)
	expiresIn := int(time.Until(expiresAt).Seconds())

	// Return tokens (OAuth2-style response)
	resp := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    expiresIn,
		Username:     user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// getClientIP extracts the client IP address from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// Fall back to RemoteAddr
	return r.RemoteAddr
}
