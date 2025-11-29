package gateway

import (
	"context"
	"net"

	"github.com/zGate-Team/zGate-Platform/internal/conn"
	"github.com/zGate-Team/zGate-Platform/internal/protocol"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Dispatcher establishes the backend connection and manages bidirectional forwarding
type Dispatcher struct {
	database   store.Database
	handler    protocol.Handler
	clientConn net.Conn
	metadata   *ConnectionMetadata
}

// NewDispatcher creates a new dispatcher for a client connection
func NewDispatcher(
	database store.Database,
	handler protocol.Handler,
	clientConn net.Conn,
	metadata *ConnectionMetadata,
) *Dispatcher {
	return &Dispatcher{
		database:   database,
		handler:    handler,
		clientConn: clientConn,
		metadata:   metadata,
	}
}

// Dispatch connects to the backend and starts proxying
func (d *Dispatcher) Dispatch(ctx context.Context) {
	// Connect to backend database
	serverConn, err := d.connectToBackend(ctx)
	if err != nil {
		utils.Logger.Error("failed to connect to backend",
			"database", d.database.Name,
			"backend_addr", d.database.BackendAddr,
			"error", err,
		)
		return
	}
	defer serverConn.Close()

	utils.Logger.Info("backend connection established",
		"database", d.database.Name,
		"backend_addr", d.database.BackendAddr,
		"client", d.metadata.ClientAddr,
	)

	// Start bidirectional forwarding
	if err := conn.Forward(d.clientConn, serverConn); err != nil {
		utils.Logger.Error("connection forwarding failed",
			"database", d.database.Name,
			"client", d.metadata.ClientAddr,
			"error", err,
		)
		return
	}

	utils.Logger.Info("connection closed",
		"database", d.database.Name,
		"client", d.metadata.ClientAddr,
	)
}

// connectToBackend establishes a connection to the backend database
func (d *Dispatcher) connectToBackend(ctx context.Context) (net.Conn, error) {
	type connResult struct {
		conn net.Conn
		err  error
	}
	connCh := make(chan connResult, 1)

	// Connect in a goroutine to allow context cancellation
	go func() {
		conn, err := d.handler.Connect(ctx, d.database.BackendAddr)
		connCh <- connResult{conn: conn, err: err}
	}()

	// Wait for connection or context cancellation
	select {
	case <-ctx.Done():
		utils.Logger.Info("backend connection cancelled",
			"database", d.database.Name,
			"reason", ctx.Err(),
		)
		return nil, ctx.Err()
	case result := <-connCh:
		return result.conn, result.err
	}
}
