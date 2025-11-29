package store

import (
	"encoding/json"
	"fmt"
)

// SaveDatabase inserts or updates a database definition.
func (s *Store) SaveDatabase(dbDef *Database) error {
	if dbDef == nil {
		return fmt.Errorf("database definition is nil")
	}

	permsJSON, err := json.Marshal(dbDef.AvailablePermissions)
	if err != nil {
		return fmt.Errorf("serialize permissions: %w", err)
	}

	encryptedPassword, err := s.encrypt([]byte(dbDef.AdminPassword))
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	query := `
	INSERT INTO databases (name, type, description, backend_addr, admin_username, admin_password, available_permissions, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT(name) DO UPDATE SET
		type=excluded.type,
		description=excluded.description,
		backend_addr=excluded.backend_addr,
		admin_username=excluded.admin_username,
		admin_password=excluded.admin_password,
		available_permissions=excluded.available_permissions,
		updated_at=CURRENT_TIMESTAMP;
	`

	if _, err := s.db.Exec(query,
		dbDef.Name,
		dbDef.Type,
		dbDef.Description,
		dbDef.BackendAddr,
		dbDef.AdminUsername,
		encryptedPassword,
		string(permsJSON),
	); err != nil {
		return fmt.Errorf("upsert database: %w", err)
	}

	return nil
}

// GetDatabase returns a single database definition.
func (s *Store) GetDatabase(name string) (*Database, error) {
	query := `SELECT name, type, description, backend_addr, admin_username, admin_password, available_permissions, created_at, updated_at FROM databases WHERE name = ?`
	var db Database
	var encrypted []byte
	var permsJSON string

	row := s.db.QueryRow(query, name)
	if err := row.Scan(&db.Name, &db.Type, &db.Description, &db.BackendAddr, &db.AdminUsername, &encrypted, &permsJSON, &db.CreatedAt, &db.UpdatedAt); err != nil {
		return nil, fmt.Errorf("fetch database: %w", err)
	}

	perms := []string{}
	if err := json.Unmarshal([]byte(permsJSON), &perms); err != nil {
		return nil, fmt.Errorf("parse permissions: %w", err)
	}
	db.AvailablePermissions = perms

	plain, err := s.decrypt(encrypted)
	if err != nil {
		return nil, fmt.Errorf("decrypt password: %w", err)
	}
	db.AdminPassword = string(plain)

	return &db, nil
}

// ListDatabases returns all defined databases.
func (s *Store) ListDatabases() ([]Database, error) {
	rows, err := s.db.Query(`SELECT name, type, description, backend_addr, admin_username, admin_password, available_permissions, created_at, updated_at FROM databases ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list databases: %w", err)
	}
	defer rows.Close()

	var result []Database
	for rows.Next() {
		var db Database
		var encrypted []byte
		var permsJSON string

		if err := rows.Scan(&db.Name, &db.Type, &db.Description, &db.BackendAddr, &db.AdminUsername, &encrypted, &permsJSON, &db.CreatedAt, &db.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan database: %w", err)
		}

		if err := json.Unmarshal([]byte(permsJSON), &db.AvailablePermissions); err != nil {
			return nil, fmt.Errorf("parse permissions: %w", err)
		}

		plain, err := s.decrypt(encrypted)
		if err != nil {
			return nil, fmt.Errorf("decrypt password: %w", err)
		}
		db.AdminPassword = string(plain)

		result = append(result, db)
	}

	return result, nil
}

// DeleteDatabase removes a database definition by name.
func (s *Store) DeleteDatabase(name string) error {
	if _, err := s.db.Exec(`DELETE FROM databases WHERE name = ?`, name); err != nil {
		return fmt.Errorf("delete database: %w", err)
	}
	return nil
}

// ListDatabaseTypes returns the distinct database types currently defined.
// Example: ["mssql", "mysql"]
func (s *Store) ListDatabaseTypes() ([]string, error) {
	rows, err := s.db.Query(`SELECT DISTINCT type FROM databases ORDER BY type`)
	if err != nil {
		return nil, fmt.Errorf("list database types: %w", err)
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("scan type: %w", err)
		}
		types = append(types, t)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate types: %w", err)
	}
	return types, nil
}
