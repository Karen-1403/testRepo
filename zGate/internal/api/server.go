package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/gateway"
	"github.com/zGate-Team/zGate-Platform/internal/policy"
	"github.com/zGate-Team/zGate-Platform/internal/proxy"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Server is the API server for CLI communication
type Server struct {
	addr          string
	server        *http.Server
	authenticator *auth.Authenticator
	policyEngine  *policy.Engine
	proxyManager  *proxy.Manager
	store         *store.Store
}

// NewServer creates a new API server
func NewServer(addr string, store *store.Store) (*Server, error) {
	// Initialize authenticator
	authenticator := auth.NewAuthenticator(store)

	// Initialize policy engine
	policyEngine := policy.NewEngine(store)

	// Initialize gateway server (for handlers)
	gwServer, err := gateway.NewServer(store)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gateway server: %w", err)
	}

	// Initialize proxy manager
	proxyManager := proxy.NewManager(store, gwServer)

	s := &Server{
		addr:          addr,
		authenticator: authenticator,
		policyEngine:  policyEngine,
		proxyManager:  proxyManager,
		store:         store,
	}

	// Setup routes
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/api/login", s.handleLogin).Methods("POST")
	router.HandleFunc("/api/refresh", s.handleRefresh).Methods("POST")
	router.HandleFunc("/api/logout", s.handleLogout).Methods("POST")

	// Protected routes (authentication required)
	router.HandleFunc("/api/active-logins", s.authMiddleware(s.handleListActiveLogins)).Methods("GET")
	router.HandleFunc("/api/active-logins/{id}", s.authMiddleware(s.handleRevokeActiveLogin)).Methods("DELETE")
	router.HandleFunc("/api/databases", s.authMiddleware(s.handleListDatabases)).Methods("GET")
	router.HandleFunc("/api/connect", s.authMiddleware(s.handleConnect)).Methods("POST")
	router.HandleFunc("/api/disconnect", s.authMiddleware(s.handleDisconnect)).Methods("POST")

	// Admin routes
	router.HandleFunc("/api/admin/login", s.handleAdminLogin).Methods("GET")
	// TODO: Implement admin authentication and uncomment these routes
	// router.HandleFunc("/api/admin/users", s.authMiddleware(s.handleListUsers)).Methods("GET")
	// router.HandleFunc("/api/admin/users/{id}", s.authMiddleware(s.handleUpdateUser)).Methods("PUT")
	// router.HandleFunc("/api/admin/users/{id}", s.authMiddleware(s.handleCreateUser)).Methods("POST")
	// router.HandleFunc("/api/admin/users/{id}", s.authMiddleware(s.handleRevokeUser)).Methods("DELETE")
	// router.HandleFunc("/api/admin/roles", s.authMiddleware(s.handleListRoles)).Methods("GET")
	// router.HandleFunc("/api/admin/roles/{id}", s.authMiddleware(s.handleUpdateRole)).Methods("PUT")
	// router.HandleFunc("/api/admin/roles/{id}", s.authMiddleware(s.handleCreateRole)).Methods("POST")
	// router.HandleFunc("/api/admin/roles/{id}", s.authMiddleware(s.handleRevokeRole)).Methods("DELETE")
	// router.HandleFunc("/api/admin/databases", s.authMiddleware(s.handleListDatabases)).Methods("GET")
	// router.HandleFunc("/api/admin/databases/{id}", s.authMiddleware(s.handleUpdateDatabase)).Methods("PUT")
	// router.HandleFunc("/api/admin/databases/{id}", s.authMiddleware(s.handleCreateDatabase)).Methods("POST")
	// router.HandleFunc("/api/admin/databases/{id}", s.authMiddleware(s.handleRevokeDatabase)).Methods("DELETE")
	// router.HandleFunc("/api/admin/active-logins", s.authMiddleware(s.handleListActiveLogins)).Methods("GET")
	// router.HandleFunc("/api/admin/active-logins/{id}", s.authMiddleware(s.handleRevokeActiveLogin)).Methods("DELETE")

	s.server = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s, nil
}

// Start starts the API server
func (s *Server) Start() error {
	utils.Logger.Info("API server starting", "addr", s.addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("API server failed: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the API server
func (s *Server) Shutdown(ctx context.Context) error {
	utils.Logger.Info("shutting down API server")
	return s.server.Shutdown(ctx)
}
