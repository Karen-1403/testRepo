package api

import (
	"encoding/json"
	"net/http"

	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// handleListDatabases handles GET /api/databases
func (s *Server) handleListDatabases(w http.ResponseWriter, r *http.Request) {
	// Get claims from context (set by authMiddleware)
	claims := r.Context().Value("claims").(*auth.Claims)

	utils.Logger.Info("listing databases", "username", claims.Username)

	// Get allowed databases from policy engine
	databases := s.policyEngine.GetAllowedDatabases(claims)

	utils.Logger.Info("databases listed", "username", claims.Username, "count", len(databases))

	// Return list
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(databases)
}