package proxy

import (
	"context"

	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/protocol"
)

// Session represents an active user session with temp database user
type Session struct {
	Username        string
	DatabaseName    string
	Port            int
	Claims          *auth.Claims
	Cancel          context.CancelFunc
	TempCredentials *protocol.TempCredentials
	DBManager       protocol.Manager
}