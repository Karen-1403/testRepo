package gateway

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/zGate-Team/zGate-Platform/internal/protocol"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Listener manages the TCP listener lifecycle for a single backend
// Used for dynamic ports created per user session
type Listener struct {
	database store.Database
	handler  protocol.Handler
	wg       sync.WaitGroup
}

// NewListener creates a new listener for the given database
func NewListener(database store.Database, handler protocol.Handler) *Listener {
	return &Listener{
		database: database,
		handler:  handler,
	}
}

// Start begins accepting connections for this backend
// Blocks until the context is cancelled
func (l *Listener) Start(ctx context.Context, listenAddr string) error {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return fmt.Errorf("failed to start listener on %s: %w", listenAddr, err)
	}

	utils.Logger.Info("dynamic listener started",
		"database", l.database.Name,
		"type", l.database.Type,
		"listen_addr", listenAddr,
	)

	// Goroutine to handle graceful shutdown
	go func() {
		<-ctx.Done()
		utils.Logger.Info("listener shutting down", "database", l.database.Name)
		listener.Close()
	}()

	// Accept connections until shutdown
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			// Check if this is due to shutdown
			select {
			case <-ctx.Done():
				utils.Logger.Info("listener closed, waiting for connections to finish",
					"database", l.database.Name,
				)
				l.wg.Wait()
				utils.Logger.Info("all connections finished", "database", l.database.Name)
				return nil
			default:
				utils.Logger.Error("failed to accept connection",
					"database", l.database.Name,
					"error", err,
				)
				continue
			}
		}

		// Create an acceptor to validate and process this connection
		acceptor := NewAcceptor(l.database, l.handler, clientConn)

		// Track active connection
		l.wg.Add(1)
		go func() {
			defer l.wg.Done()
			acceptor.Accept(ctx)
		}()
	}
}
