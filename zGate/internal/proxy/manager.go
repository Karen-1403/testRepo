package proxy

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/gateway"
	"github.com/zGate-Team/zGate-Platform/internal/protocol"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Manager manages dynamic proxy sessions
type Manager struct {
	sessions map[string]*Session
	store    *store.Store
	gwServer *gateway.Server
	mu       sync.RWMutex
}

// NewManager creates a new proxy manager
func NewManager(store *store.Store, gwServer *gateway.Server) *Manager {
	return &Manager{
		sessions: make(map[string]*Session),
		store:    store,
		gwServer: gwServer,
	}
}

// StartSession creates a new dynamic proxy with temp database user
func (m *Manager) StartSession(token string, claims *auth.Claims, databaseName string) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if session already exists
	if existing, exists := m.sessions[token]; exists {
		return existing, nil
	}

	// Find database config
	database, err := m.store.GetDatabase(databaseName)
	if err != nil {
		return nil, fmt.Errorf("database not found: %s", databaseName)
	}

	// Create DB manager
	dbMgr, err := protocol.NewManager(*database)
	if err != nil {
		return nil, fmt.Errorf("failed to create DB manager: %w", err)
	}

	// Generate temp credentials
	baseUsername := extractUsername(claims.Username)
	tempUsername := protocol.GenerateTempUsername(baseUsername)
	tempPassword := protocol.GenerateTempPassword()

	tempCreds := &protocol.TempCredentials{
		Username: tempUsername,
		Password: tempPassword,
	}

	utils.Logger.Info("generated temp credentials",
		"database", databaseName,
		"zgate_user", claims.Username,
		"temp_user", tempUsername,
	)

	// Fetch user permissions from store (real-time, not from cached claims)
	user, err := m.store.GetUser(claims.Username)
	if err != nil {
		dbMgr.Close()
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get role permissions
	rolePerms, err := m.store.GetPermissionsForRoles(user.Roles)
	if err != nil {
		dbMgr.Close()
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	// Combine role permissions and custom permissions
	allPermissions := append(rolePerms, user.CustomPermissions...)

	// Determine permissions for this specific database
	permissions := determinePermissionsForDatabase(allPermissions, databaseName)

	// Create temp user in database
	ctx := context.Background()
	if err := dbMgr.CreateTempUser(ctx, tempUsername, tempPassword, permissions); err != nil {
		dbMgr.Close()
		return nil, fmt.Errorf("failed to create temp user: %w", err)
	}

	// Find available port
	port, err := getFreePort()
	if err != nil {
		dbMgr.DeleteTempUser(ctx, tempUsername)
		dbMgr.Close()
		return nil, fmt.Errorf("failed to find free port: %w", err)
	}

	// Create session
	ctx, cancel := context.WithCancel(context.Background())
	session := &Session{
		Username:        claims.Username,
		DatabaseName:    databaseName,
		Port:            port,
		Claims:          claims,
		Cancel:          cancel,
		TempCredentials: tempCreds,
		DBManager:       dbMgr,
	}

	// Start dynamic proxy in background
	go m.startDynamicProxy(ctx, session, database)

	// Store session
	m.sessions[token] = session

	utils.Logger.Info("session started",
		"zgate_user", claims.Username,
		"database", databaseName,
		"port", port,
		"temp_user", tempUsername,
	)

	return session, nil
}

// StopSession stops a session and deletes temp database user
func (m *Manager) StopSession(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[token]
	if !exists {
		return fmt.Errorf("session not found")
	}

	utils.Logger.Info("stopping session",
		"zgate_user", session.Username,
		"database", session.DatabaseName,
		"temp_user", session.TempCredentials.Username,
	)

	// Delete temp database user
	ctx := context.Background()
	if err := session.DBManager.DeleteTempUser(ctx, session.TempCredentials.Username); err != nil {
		utils.Logger.Error("failed to delete temp user", "error", err)
	}

	// Close DB manager connection
	if err := session.DBManager.Close(); err != nil {
		utils.Logger.Error("failed to close DB manager", "error", err)
	}

	// Cancel context (stops the listener)
	session.Cancel()

	// Remove from map
	delete(m.sessions, token)

	utils.Logger.Info("session stopped",
		"zgate_user", session.Username,
		"database", session.DatabaseName,
	)

	return nil
}

// startDynamicProxy starts a listener on the dynamic port
func (m *Manager) startDynamicProxy(ctx context.Context, session *Session, database *store.Database) {
	listenAddr := fmt.Sprintf(":%d", session.Port)

	handler := m.gwServer.GetHandler(database.Type)
	if handler == nil {
		utils.Logger.Error("handler not found", "type", database.Type)
		return
	}

	listener := gateway.NewListener(*database, handler)
	if err := listener.Start(ctx, listenAddr); err != nil {
		utils.Logger.Error("dynamic proxy stopped", "error", err)
	}
}

// Helper functions

func extractUsername(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		return parts[0]
	}
	return email
}

func determinePermissionsForDatabase(permissions []store.Permission, databaseName string) []string {
	for _, perm := range permissions {
		if perm.Database == databaseName {
			return []string{perm.Level}
		}
	}
	return []string{}
}

func getFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	return addr.Port, nil
}
