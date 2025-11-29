package api

import (
	"encoding/json"
	"net/http"

	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// ConnectRequest represents connect request payload
type ConnectRequest struct {
	DatabaseName string `json:"database_name"`
}

// ConnectResponse represents connect response payload
type ConnectResponse struct {
	Port         int    `json:"port"`
	DatabaseName string `json:"database_name"`
	Message      string `json:"message"`
	TempUsername string `json:"temp_username"`
	TempPassword string `json:"temp_password"`
}

// handleConnect handles POST /api/connect
func (s *Server) handleConnect(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	token := r.Context().Value("token").(string)

	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	utils.Logger.Info("connect request", "username", claims.Username, "database", req.DatabaseName)

	// Check permission
	if !s.policyEngine.CanAccess(claims, req.DatabaseName) {
		utils.Logger.Warn("access denied", "username", claims.Username, "database", req.DatabaseName)
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	// Start proxy session (creates temp DB user)
	session, err := s.proxyManager.StartSession(token, claims, req.DatabaseName)
	if err != nil {
		utils.Logger.Error("failed to start session", "error", err)
		http.Error(w, "failed to start proxy", http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("proxy session started",
		"username", claims.Username,
		"database", req.DatabaseName,
		"port", session.Port,
		"temp_user", session.TempCredentials.Username,
	)

	// Return connection info with temp credentials
	resp := ConnectResponse{
		Port:         session.Port,
		DatabaseName: req.DatabaseName,
		Message:      "Proxy started successfully",
		TempUsername: session.TempCredentials.Username,
		TempPassword: session.TempCredentials.Password,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// DisconnectRequest represents disconnect request payload
type DisconnectRequest struct {
	DatabaseName string `json:"database_name"`
}

// handleDisconnect handles POST /api/disconnect
func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*auth.Claims)
	token := r.Context().Value("token").(string)

	var req DisconnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	utils.Logger.Info("disconnect request", "username", claims.Username, "database", req.DatabaseName)

	// Stop session (deletes temp user)
	if err := s.proxyManager.StopSession(token); err != nil {
		utils.Logger.Error("failed to stop session", "error", err)
		http.Error(w, "failed to stop proxy", http.StatusInternalServerError)
		return
	}

	utils.Logger.Info("proxy session stopped", "username", claims.Username, "database", req.DatabaseName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Disconnected successfully",
	})
}