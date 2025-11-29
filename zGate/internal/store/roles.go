package store

import (
	"context"
	"fmt"
	"strings"
)

// SaveRole inserts or updates a role and its permissions.
func (s *Store) SaveRole(role *Role) error {
	if role == nil {
		return fmt.Errorf("role is nil")
	}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(`
		INSERT INTO roles (name, description) VALUES (?, ?)
		ON CONFLICT(name) DO UPDATE SET description=excluded.description
	`, role.Name, role.Description); err != nil {
		return err
	}

	if _, err = tx.Exec(`DELETE FROM role_permissions WHERE role_name = ?`, role.Name); err != nil {
		return err
	}

	for _, perm := range role.Permissions {
		if _, err = tx.Exec(`INSERT INTO role_permissions (role_name, database_name, level) VALUES (?, ?, ?)`, role.Name, perm.Database, perm.Level); err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

// GetRole fetches a role by name.
func (s *Store) GetRole(name string) (*Role, error) {
	row := s.db.QueryRow(`SELECT name, description FROM roles WHERE name = ?`, name)
	var role Role
	if err := row.Scan(&role.Name, &role.Description); err != nil {
		return nil, fmt.Errorf("fetch role: %w", err)
	}

	perms, err := s.getPermissionsForRole(name)
	if err != nil {
		return nil, err
	}
	role.Permissions = perms

	return &role, nil
}

// ListRoles returns all roles with permissions.
func (s *Store) ListRoles() ([]Role, error) {
	rows, err := s.db.Query(`SELECT name, description FROM roles ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.Name, &role.Description); err != nil {
			return nil, err
		}

		perms, err := s.getPermissionsForRole(role.Name)
		if err != nil {
			return nil, err
		}
		role.Permissions = perms

		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Store) getPermissionsForRole(roleName string) ([]Permission, error) {
	rows, err := s.db.Query(`SELECT database_name, level FROM role_permissions WHERE role_name = ?`, roleName)
	if err != nil {
		return nil, fmt.Errorf("role permissions: %w", err)
	}
	defer rows.Close()

	var perms []Permission
	for rows.Next() {
		var perm Permission
		if err := rows.Scan(&perm.Database, &perm.Level); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	return perms, nil
}

// GetPermissionsForRoles aggregates permissions for the supplied roles.
func (s *Store) GetPermissionsForRoles(roleNames []string) ([]Permission, error) {
	if len(roleNames) == 0 {
		return nil, nil
	}

	placeholders := strings.Repeat("?,", len(roleNames))
	placeholders = strings.TrimSuffix(placeholders, ",")

	query := fmt.Sprintf(`SELECT database_name, level FROM role_permissions WHERE role_name IN (%s)`, placeholders)
	args := make([]any, len(roleNames))
	for i, name := range roleNames {
		args[i] = name
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("aggregate perms: %w", err)
	}
	defer rows.Close()

	var perms []Permission
	for rows.Next() {
		var perm Permission
		if err := rows.Scan(&perm.Database, &perm.Level); err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	return perms, nil
}

// GetUsersForRole lists all usernames that are directly assigned the specified role.
// It does not include users who might receive equivalent permissions through other roles.
func (s *Store) GetUsersForRole(roleName string) ([]string, error) {
	if roleName == "" {
		return nil, fmt.Errorf("role name is empty")
	}

	rows, err := s.db.Query(`SELECT username FROM user_roles WHERE role_name = ? ORDER BY username`, roleName)
	if err != nil {
		return nil, fmt.Errorf("list users for role: %w", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, fmt.Errorf("scan username: %w", err)
		}
		users = append(users, username)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate usernames: %w", err)
	}
	return users, nil
}
