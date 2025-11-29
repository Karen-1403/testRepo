package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Manager implements protocol.Manager for MySQL
type Manager struct {
	database store.Database
	db       *sql.DB
}

// NewManager creates a new MySQL manager
func NewManager(database store.Database) (*Manager, error) {
	connString := fmt.Sprintf(
		"%s:%s@tcp(%s)/mysql",
		database.AdminUsername,
		database.AdminPassword,
		database.BackendAddr,
	)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	utils.Logger.Info("MySQL manager connected", "database", database.Name)

	return &Manager{
		database: database,
		db:       db,
	}, nil
}

// CreateTempUser creates a temporary MySQL user
func (m *Manager) CreateTempUser(ctx context.Context, username, password string, permissions []string) error {
	utils.Logger.Info("creating temp MySQL user", "database", m.database.Name, "username", username)

	// Create USER
	createUserSQL := fmt.Sprintf(
		"CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'",
		username, password,
	)

	if _, err := m.db.ExecContext(ctx, createUserSQL); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Grant permissions
	for _, perm := range permissions {
		var grantSQL string

		switch perm {
		case "read":
			grantSQL = fmt.Sprintf("GRANT SELECT ON *.* TO '%s'@'%%'", username)
		case "write":
			grantSQL = fmt.Sprintf("GRANT SELECT, INSERT, UPDATE, DELETE ON *.* TO '%s'@'%%'", username)
		case "admin":
			grantSQL = fmt.Sprintf("GRANT ALL PRIVILEGES ON *.* TO '%s'@'%%'", username)
		default:
			utils.Logger.Warn("unknown permission", "permission", perm)
			continue
		}

		if _, err := m.db.ExecContext(ctx, grantSQL); err != nil {
			utils.Logger.Error("failed to grant permission", "permission", perm, "error", err)
		}
	}

	m.db.ExecContext(ctx, "FLUSH PRIVILEGES")

	utils.Logger.Info("temp MySQL user created", "database", m.database.Name, "username", username)
	return nil
}

// DeleteTempUser removes a temporary MySQL user
func (m *Manager) DeleteTempUser(ctx context.Context, username string) error {
	utils.Logger.Info("deleting temp MySQL user", "database", m.database.Name, "username", username)
	dropUserSQL := fmt.Sprintf("DROP USER IF EXISTS '%s'@'%%'", username)

	if _, err := m.db.ExecContext(ctx, dropUserSQL); err != nil {
		return fmt.Errorf("failed to drop user: %w", err)
	}

	m.db.ExecContext(ctx, "FLUSH PRIVILEGES")

	utils.Logger.Info("temp MySQL user deleted", "database", m.database.Name, "username", username)
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
	return "mysql"
}
