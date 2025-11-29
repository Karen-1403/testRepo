package gateway

import (
	"fmt"
	"sync"

	"github.com/zGate-Team/zGate-Platform/internal/protocol"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Server orchestrates dynamic backend listeners
// Note: No static listeners anymore - all are created on-demand per session
type Server struct {
	handlers map[string]protocol.Handler
	mu       sync.RWMutex
}

// NewServer creates a new gateway server instance
// Pre-initializes handlers for all database types
func NewServer(store *store.Store) (*Server, error) {
	if store == nil {
		return nil, fmt.Errorf("store cannot be nil")
	}

	s := &Server{
		handlers: make(map[string]protocol.Handler),
	}

	// Pre-initialize handlers for each unique database type
	seenTypes := make(map[string]bool)
	databases, err := store.ListDatabases()
	if err != nil {
		return nil, fmt.Errorf("failed to list databases: %w", err)
	}

	for _, db := range databases {
		if seenTypes[db.Type] {
			continue
		}
		seenTypes[db.Type] = true

		h, err := protocol.NewHandler(db.Type)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize handler for type %s: %w", db.Type, err)
		}
		s.handlers[db.Type] = h
		utils.Logger.Info("handler initialized", "type", db.Type)
	}

	return s, nil
}

// GetHandler retrieves the pre-initialized handler for a database type
func (s *Server) GetHandler(dbType string) protocol.Handler {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.handlers[dbType]
}
