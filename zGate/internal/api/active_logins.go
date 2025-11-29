package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// ActiveLoginInfo represents session information for display
type ActiveLoginInfo struct {
	ID         int64  `json:"id"`
	CreatedAt  string `json:"created_at"`
	LastUsedAt string `json:"last_used_at"`
	ExpiresAt  string `json:"expires_at"`
	UserAgent  string `json:"user_agent,omitempty"`
	IPAddress  string `json:"ip_address,omitempty"`
	IsCurrent  bool   `json:"is_current"`
}

// ActiveLoginsResponse represents the list of active logins
type ActiveLoginsResponse struct {
	ActiveLogins []ActiveLoginInfo `json:"active_logins"`
	Total        int               `json:"total"`
}

// RevokeActiveLoginResponse represents session revocation response
type RevokeActiveLoginResponse struct {
	Message string `json:"message"`
}

// handleListActiveLogins handles GET /api/active-logins
func (s *Server) handleListActiveLogins(w http.ResponseWriter, r *http.Request) {
	// Get username from context (set by auth middleware)
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all active sessions for the user
	sessions, err := s.store.GetUserActiveLogins(claims.Username)
	if err != nil {
		utils.Logger.Error("failed to get user active logins", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Get current token from Authorization header to mark current session
	currentToken := extractTokenFromHeader(r)

	// Convert to response format
	loginInfos := make([]ActiveLoginInfo, 0, len(sessions))
	for _, session := range sessions {
		// Check if this is the current session (simplified - in production, compare token hashes)
		isCurrent := false
		if currentToken != "" {
			// Note: This is a simplified check. In production, you'd validate the token
			// and compare its hash with session.TokenHash
			isCurrent = false // For now, we'll mark based on most recent usage
		}

		loginInfos = append(loginInfos, ActiveLoginInfo{
			ID:         session.ID,
			CreatedAt:  session.CreatedAt.Format("2006-01-02 15:04:05"),
			LastUsedAt: session.LastUsedAt.Format("2006-01-02 15:04:05"),
			ExpiresAt:  session.ExpiresAt.Format("2006-01-02 15:04:05"),
			UserAgent:  session.UserAgent,
			IPAddress:  session.IPAddress,
			IsCurrent:  isCurrent,
		})
	}

	resp := ActiveLoginsResponse{
		ActiveLogins: loginInfos,
		Total:        len(loginInfos),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// handleRevokeActiveLogin handles DELETE /api/active-logins/{id}
func (s *Server) handleRevokeActiveLogin(w http.ResponseWriter, r *http.Request) {
	// Get username from context (set by auth middleware)
	claims, ok := r.Context().Value("claims").(*auth.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session ID from URL
	vars := mux.Vars(r)
	sessionIDStr := vars["id"]
	sessionID, err := strconv.ParseInt(sessionIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid session ID", http.StatusBadRequest)
		return
	}

	// Revoke the session (ensures user can only revoke their own sessions)
	if err := s.store.RevokeActiveLoginByID(sessionID, claims.Username); err != nil {
		utils.Logger.Warn("failed to revoke active login", "error", err, "session_id", sessionID)
		http.Error(w, "session not found or already revoked", http.StatusNotFound)
		return
	}

	utils.Logger.Info("active login revoked", "username", claims.Username, "session_id", sessionID)

	resp := RevokeActiveLoginResponse{
		Message: "Active login revoked successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// extractTokenFromHeader extracts the bearer token from Authorization header
func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
