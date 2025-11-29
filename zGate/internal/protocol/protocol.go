package protocol

import (
	"context"
	"fmt"
	"net"

	"github.com/zGate-Team/zGate-Platform/internal/store"
)

// ---------------------------------------------------------
// 1. The Core Interface (Contract)
// ---------------------------------------------------------

// DatabaseHandler defines the unified interface for database wire protocol interactions.
// It exposes granular methods so the Dispatcher can orchestrate the connection lifecycle.
// Every supported database (MySQL, MSSQL, Postgres) must implement this interface.
type DatabaseHandler interface {
	// --- Connection Lifecycle ---

	// Connect establishes a TCP connection to the backend database.
	Connect(ctx context.Context, addr string) (net.Conn, error)
	// ConnectWithCredentials establishes a TCP connection to the backend database with credentials.
	ConnectWithCredentials(ctx context.Context, addr, username, password string) (net.Conn, error)

	// Handshake performs the initial authentication handshake on both sides.
	// It returns context about the connection (e.g., "Handshake Complete").
	Handshake(ctx context.Context, clientConn, serverConn net.Conn) error

	// --- Command Loop Primitives ---

	// ReadCommand reads a command packet from the client.
	// It returns the parsed query string (if it is a query) or an empty string (if it's another command like Quit/Ping).
	// It also returns the raw packet bytes for forwarding.
	ReadCommand(clientConn net.Conn) (query string, packet []byte, err error)

	// SendError sends a protocol-specific error message to the client.
	// Used when the Dispatcher blocks a query.
	SendError(clientConn net.Conn, errMsg string) error

	// ForwardCommand sends the raw command packet to the server.
	ForwardCommand(serverConn net.Conn, packet []byte) error

	// ForwardResult reads the response from the server and sends it to the client.
	// This is where Data Masking logic will live inside the implementation.
	ForwardResult(clientConn, serverConn net.Conn) error

	// --- User Management ---

	// CreateTempUser creates a short-lived user for a specific session.
	CreateTempUser(ctx context.Context, username, password string, permissions []string) error

	// DeleteTempUser removes the temporary database user.
	DeleteTempUser(ctx context.Context, username string) error

	// --- Utilities ---

	// GetType returns the database type string (e.g., "mysql").
	GetType() string

	// Close cleans up any shared resources (e.g., admin connections).
	Close() error
}

// ---------------------------------------------------------
// 2. The Parser Helper Interface
// ---------------------------------------------------------

// Parser defines the interface for inspecting protocol packets.
// Implementations use this internally to extract SQL from raw bytes.
type Parser interface {
	// ParseQuery attempts to extract a SQL string from a raw packet.
	ParseQuery(packet []byte) (string, bool)
}

// ---------------------------------------------------------
// 3. The Factory Function (Creator)
// ---------------------------------------------------------

// NewDatabaseHandler creates a unified handler for the specified database configuration.
// This replaces the separate NewManager and NewHandler functions.
func NewDatabaseHandler(database store.Database) (DatabaseHandler, error) {
	switch database.Type {
	case "mysql":
		// return mysql.NewHandler(database) // TODO: Uncomment when implemented
		return nil, nil
	case "mssql":
		// return mssql.NewHandler(database) // TODO: Uncomment when implemented
		return nil, nil
	// case "postgres":
	// 	return postgres.NewHandler(database), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", database.Type)
	}
}
