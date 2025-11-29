package mssql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Manager implements protocol.Manager for MSSQL
type Manager struct {
	database store.Database
	db       *sql.DB
}

// NewManager creates a new MSSQL manager
func NewManager(database store.Database) (*Manager, error) {
	connString := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;database=master",
		database.BackendAddr,
		database.AdminUsername,
		database.AdminPassword,
	)

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MSSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MSSQL: %w", err)
	}

	utils.Logger.Info("MSSQL manager connected", "database", database.Name)

	return &Manager{
		database: database,
		db:       db,
	}, nil
}

// CreateTempUser creates a temporary MSSQL login and user
func (m *Manager) CreateTempUser(ctx context.Context, username, password string, permissions []string) error {
	utils.Logger.Info("creating temp MSSQL user",
		"database", m.database.Name,
		"username", username,
	)

	// Create LOGIN
	createLoginSQL := fmt.Sprintf(
		`IF NOT EXISTS (SELECT * FROM sys.server_principals WHERE name = '%s')
		BEGIN
			CREATE LOGIN [%s] WITH PASSWORD = '%s'
		END`,
		username, username, password,
	)

	if _, err := m.db.ExecContext(ctx, createLoginSQL); err != nil {
		return fmt.Errorf("failed to create login: %w", err)
	}

	// Create USER
	createUserSQL := fmt.Sprintf(
		`IF NOT EXISTS (SELECT * FROM sys.database_principals WHERE name = '%s')
		BEGIN
			CREATE USER [%s] FOR LOGIN [%s]
		END`,
		username, username, username,
	)

	if _, err := m.db.ExecContext(ctx, createUserSQL); err != nil {
		m.DeleteTempUser(ctx, username)
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Grant permissions
	for _, perm := range permissions {
		var grantSQL string

		switch perm {
		case "read":
			grantSQL = fmt.Sprintf(`GRANT SELECT TO [%s]`, username)
		case "write":
			grantSQL = fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE TO [%s]`, username)
		case "admin":
			grantSQL = fmt.Sprintf(`ALTER ROLE db_owner ADD MEMBER [%s]`, username)
		default:
			utils.Logger.Warn("unknown permission", "permission", perm)
			continue
		}

		if _, err := m.db.ExecContext(ctx, grantSQL); err != nil {
			utils.Logger.Error("failed to grant permission", "permission", perm, "error", err)
		}
	}

	utils.Logger.Info("temp MSSQL user created", "database", m.database.Name, "username", username)
	return nil
}

// DeleteTempUser removes a temporary MSSQL user
func (m *Manager) DeleteTempUser(ctx context.Context, username string) error {
	utils.Logger.Info("deleting temp MSSQL user", "database", m.database.Name, "username", username)

	// Drop USER
	dropUserSQL := fmt.Sprintf(
		`IF EXISTS (SELECT * FROM sys.database_principals WHERE name = '%s')
		BEGIN
			DROP USER [%s]
		END`,
		username, username,
	)
	m.db.ExecContext(ctx, dropUserSQL)

	// Drop LOGIN
	dropLoginSQL := fmt.Sprintf(
		`IF EXISTS (SELECT * FROM sys.server_principals WHERE name = '%s')
		BEGIN
			DROP LOGIN [%s]
		END`,
		username, username,
	)

	if _, err := m.db.ExecContext(ctx, dropLoginSQL); err != nil {
		return fmt.Errorf("failed to drop login: %w", err)
	}

	utils.Logger.Info("temp MSSQL user deleted", "database", m.database.Name, "username", username)
	return nil
}

// Close closes the database connection
func (m *Manager) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

// GetType returns the database type
func (m *Manager) GetType() string {
	return "mssql"
}
