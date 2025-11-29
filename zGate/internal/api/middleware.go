package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// authMiddleware validates JWT token and adds claims to context
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.Logger.Warn("missing authorization header")
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Check Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Logger.Warn("invalid authorization header format")
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := auth.ValidateToken(token)
		if err != nil {
			utils.Logger.Warn("invalid token", "error", err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		utils.Logger.Info("token validated", "username", claims.Username)

		// Add claims and token to context
		ctx := context.WithValue(r.Context(), "claims", claims)
		ctx = context.WithValue(ctx, "token", token)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}