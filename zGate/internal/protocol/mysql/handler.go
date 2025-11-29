package mysql

import (
	"context"
	"fmt"
	"net"
)

// Handler implements protocol.Handler for MySQL
type Handler struct{}

// NewHandler creates a new MySQL handler
func NewHandler() *Handler {
	return &Handler{}
}

// Connect establishes a TCP connection to MySQL server
func (h *Handler) Connect(ctx context.Context, addr string) (net.Conn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	return conn, nil
}

// GetType returns the database type
func (h *Handler) GetType() string {
	return "mysql"
}

// Close closes any resources
func (h *Handler) Close() error {
	return nil
}