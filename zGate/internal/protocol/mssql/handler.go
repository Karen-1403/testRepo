package mssql

import (
	"context"
	"fmt"
	"net"
)

// Handler implements protocol.Handler for MSSQL
type Handler struct{}

// NewHandler creates a new MSSQL handler
func NewHandler() *Handler {
	return &Handler{}
}

// Connect establishes a TCP connection to MSSQL server
func (h *Handler) Connect(ctx context.Context, addr string) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MSSQL: %w", err)
	}
	return conn, nil
}

// GetType returns the database type
func (h *Handler) GetType() string {
	return "mssql"
}

// Close closes any resources
func (h *Handler) Close() error {
	return nil
}