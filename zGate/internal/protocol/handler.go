package protocol

import (
	"context"
	"fmt"
	"net"

	"github.com/zGate-Team/zGate-Platform/internal/protocol/mssql"
	"github.com/zGate-Team/zGate-Platform/internal/protocol/mysql"
)

// Handler defines the interface for database connection handlers
type Handler interface {
	// Connect establishes a connection to the database
	Connect(ctx context.Context, addr string) (net.Conn, error)

	// GetType returns the database type (mssql, mysql, postgres)
	GetType() string

	// Close closes any resources held by the handler
	Close() error
}

// NewHandler creates a handler for the specified database type
func NewHandler(dbType string) (Handler, error) {
	switch dbType {
	case "mssql":
		return mssql.NewHandler(), nil
	case "mysql":
		return mysql.NewHandler(), nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}