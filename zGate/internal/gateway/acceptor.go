package gateway

import (
	"context"
	"net"
	"time"

	"github.com/zGate-Team/zGate-Platform/internal/protocol"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Acceptor handles a single client connection from acceptance to completion
type Acceptor struct {
	database   store.Database
	handler    protocol.Handler
	clientConn net.Conn
}

// NewAcceptor creates a new acceptor for a client connection
func NewAcceptor(database store.Database, handler protocol.Handler, clientConn net.Conn) *Acceptor {
	return &Acceptor{
		database:   database,
		handler:    handler,
		clientConn: clientConn,
	}
}

// Accept processes the client connection
func (a *Acceptor) Accept(ctx context.Context) {
	defer a.clientConn.Close()

	clientAddr := a.clientConn.RemoteAddr().String()

	utils.Logger.Info("connection accepted",
		"database", a.database.Name,
		"client", clientAddr,
	)

	// Create connection metadata
	connMeta := &ConnectionMetadata{
		ClientAddr:   clientAddr,
		DatabaseName: a.database.Name,
		DatabaseType: a.database.Type,
		ConnectedAt:  time.Now(),
	}

	// Delegate to dispatcher for actual proxying
	dispatcher := NewDispatcher(a.database, a.handler, a.clientConn, connMeta)
	dispatcher.Dispatch(ctx)
}

// ConnectionMetadata holds metadata about a client connection
type ConnectionMetadata struct {
	ClientAddr   string
	DatabaseName string
	DatabaseType string
	ConnectedAt  time.Time
}
