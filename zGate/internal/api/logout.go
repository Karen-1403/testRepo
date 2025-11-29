package api

import (
	"encoding/json"
	"net/http"

	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// LogoutRequest represents logout request payload
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// LogoutResponse represents logout response payload
type LogoutResponse struct {
	Message string `json:"message"`
}

// handleLogout handles POST /api/logout
func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	var req LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Logger.Error("invalid logout request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		http.Error(w, "refresh_token is required", http.StatusBadRequest)
		return
	}

	// Revoke the refresh token
	if err := s.store.RevokeRefreshToken(req.RefreshToken); err != nil {
		utils.Logger.Warn("failed to revoke refresh token", "error", err)
		// Return success anyway - token might already be revoked or expired
	}

	utils.Logger.Info("user logged out")

	// Return success
	resp := LogoutResponse{
		Message: "Logged out successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
