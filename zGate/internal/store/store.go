package store

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Store owns the SQLite connection and exposes CRUD helpers for zGate metadata.
type Store struct {
	db            *sql.DB
	encryptionKey []byte
}

// NewStore opens (or creates) the SQLite database at dbPath and prepares the schema.
func NewStore(dbPath string, encryptionKey []byte) (*Store, error) {
	if len(encryptionKey) != 32 {
		return nil, errors.New("encryption key must be exactly 32 bytes for AES-256")
	}

	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, fmt.Errorf("failed creating data directory: %w", err)
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_busy_timeout=5000&_foreign_keys=on", dbPath))
	if err != nil {
		return nil, fmt.Errorf("failed opening sqlite database: %w", err)
	}

	store := &Store{db: db, encryptionKey: encryptionKey}
	if err := store.initSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

// Close releases the database resources.
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS databases (
		name TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		description TEXT,
		backend_addr TEXT NOT NULL,
		admin_username TEXT NOT NULL,
		admin_password BLOB NOT NULL,
		available_permissions TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS roles (
		name TEXT PRIMARY KEY,
		description TEXT
	);

	CREATE TABLE IF NOT EXISTS role_permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		role_name TEXT NOT NULL,
		database_name TEXT NOT NULL,
		level TEXT NOT NULL,
		UNIQUE(role_name, database_name),
		FOREIGN KEY(role_name) REFERENCES roles(name) ON DELETE CASCADE,
		FOREIGN KEY(database_name) REFERENCES databases(name) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password_hash BLOB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS admins (
		username TEXT PRIMARY KEY,
		password_hash BLOB NOT NULL,
		name TEXT NOT NULL,
		email TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_login_at TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_roles (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		role_name TEXT NOT NULL,
		UNIQUE(username, role_name),
		FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE,
		FOREIGN KEY(role_name) REFERENCES roles(name) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS user_custom_permissions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL,
		database_name TEXT NOT NULL,
		level TEXT NOT NULL,
		UNIQUE(username, database_name, level),
		FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS refresh_tokens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		token_hash TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		revoked BOOLEAN DEFAULT 0,
		revoked_at TIMESTAMP,
		user_agent TEXT,
		ip_address TEXT,
		FOREIGN KEY(username) REFERENCES users(username) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_username ON refresh_tokens(username);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
	CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
	`

	if _, err := s.db.Exec(schema); err != nil {
		return fmt.Errorf("failed applying schema: %w", err)
	}

	return nil
}
