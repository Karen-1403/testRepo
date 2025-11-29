package protocol

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/zGate-Team/zGate-Platform/internal/protocol/mssql"
	"github.com/zGate-Team/zGate-Platform/internal/protocol/mysql"
	"github.com/zGate-Team/zGate-Platform/internal/store"
)

// Manager defines the interface for database user management
type Manager interface {
	// CreateTempUser creates a temporary database user
	CreateTempUser(ctx context.Context, username, password string, permissions []string) error

	// DeleteTempUser removes a temporary database user
	DeleteTempUser(ctx context.Context, username string) error

	// Close closes the database connection
	Close() error

	// GetType returns the database type
	GetType() string
}

// TempCredentials represents temporary user credentials
type TempCredentials struct {
	Username string
	Password string
}

// NewManager creates a manager for the specified database
func NewManager(database store.Database) (Manager, error) {
	switch database.Type {
	case "mssql":
		return mssql.NewManager(database)
	case "mysql":
		return mysql.NewManager(database)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", database.Type)
	}
}

// GenerateTempUsername creates a unique temporary username
func GenerateTempUsername(baseUsername string) string {
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	suffix := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("zgate_%s_%s", baseUsername, suffix)
}

// GenerateTempPassword creates a secure random password
func GenerateTempPassword() string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	password := hex.EncodeToString(randomBytes)

	// FIX: Added 'Z' and 'G' prefixes to ensuring Uppercase characters are present
	// Format: Z<8chars>#<8chars>$G<8chars>
	// This satisfies MySQL MEDIUM policy (Length, Number, Mixed Case, Special Char)
	return fmt.Sprintf("Z%s#%s$G%s", password[:8], password[8:16], password[16:24])
}
