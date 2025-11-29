package store

import (
	"context"
	"fmt"
)

// SaveUser writes the user record (expects PasswordHash to be populated) and replaces roles/custom permissions.
func (s *Store) SaveUser(user *User) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}
	if user.PasswordHash == "" {
		return fmt.Errorf("user encrypted password is required")
	}

	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err = tx.Exec(`
		INSERT INTO users (username, password_hash)
		VALUES (?, ?)
		ON CONFLICT(username) DO UPDATE SET password_hash=excluded.password_hash
	`, user.Username, user.PasswordHash); err != nil {
		return err
	}

	if _, err = tx.Exec(`DELETE FROM user_roles WHERE username = ?`, user.Username); err != nil {
		return err
	}
	for _, role := range user.Roles {
		if _, err = tx.Exec(`INSERT INTO user_roles (username, role_name) VALUES (?, ?)`, user.Username, role); err != nil {
			return err
		}
	}

	if _, err = tx.Exec(`DELETE FROM user_custom_permissions WHERE username = ?`, user.Username); err != nil {
		return err
	}
	for _, perm := range user.CustomPermissions {
		if _, err = tx.Exec(`INSERT INTO user_custom_permissions (username, database_name, level) VALUES (?, ?, ?)`, user.Username, perm.Database, perm.Level); err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

// CreateUserWithPassword encrypts the password then delegates to SaveUser.
func (s *Store) CreateUserWithPassword(username, plainPassword string, roles []string, custom []Permission) error {
	encrypted, err := s.encrypt([]byte(plainPassword))
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	user := &User{
		Username:          username,
		PasswordHash:      string(encrypted),
		Roles:             roles,
		CustomPermissions: custom,
	}
	return s.SaveUser(user)
}

// GetUser fetches the user with roles and custom permissions.
func (s *Store) GetUser(username string) (*User, error) {
	row := s.db.QueryRow(`SELECT username, password_hash, created_at FROM users WHERE username = ?`, username)
	var user User
	if err := row.Scan(&user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("fetch user: %w", err)
	}

	roles, err := s.getUserRoles(username)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	custom, err := s.getUserCustomPermissions(username)
	if err != nil {
		return nil, err
	}
	user.CustomPermissions = custom

	return &user, nil
}

// ListUsers returns all users with their roles and custom permissions.
func (s *Store) ListUsers() ([]User, error) {
	rows, err := s.db.Query(`SELECT username, password_hash, created_at FROM users ORDER BY username`)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Username, &user.PasswordHash, &user.CreatedAt); err != nil {
			return nil, err
		}

		user.Roles, err = s.getUserRoles(user.Username)
		if err != nil {
			return nil, err
		}

		user.CustomPermissions, err = s.getUserCustomPermissions(user.Username)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// VerifyPassword decrypts and compares a user's password against the supplied password.
func (s *Store) VerifyPassword(username, plainPassword string) error {
	user, err := s.GetUser(username)
	if err != nil {
		return err
	}

	decrypted, err := s.decrypt([]byte(user.PasswordHash))
	if err != nil {
		return fmt.Errorf("decrypt password: %w", err)
	}

	if string(decrypted) != plainPassword {
		return fmt.Errorf("invalid password")
	}
	return nil
}

// SetUserPassword updates an existing user's password.
func (s *Store) SetUserPassword(username, plainPassword string) error {
	encrypted, err := s.encrypt([]byte(plainPassword))
	if err != nil {
		return fmt.Errorf("encrypt password: %w", err)
	}

	if _, err := s.db.Exec(`UPDATE users SET password_hash = ? WHERE username = ?`, encrypted, username); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return nil
}

func (s *Store) getUserRoles(username string) ([]string, error) {
	rows, err := s.db.Query(`SELECT role_name FROM user_roles WHERE username = ? ORDER BY role_name`, username)
	if err != nil {
		return nil, fmt.Errorf("user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (s *Store) getUserCustomPermissions(username string) ([]Permission, error) {
	rows, err := s.db.Query(`SELECT database_name, level FROM user_custom_permissions WHERE username = ?`, username)
	if err != nil {
		return nil, fmt.Errorf("user custom perms: %w", err)
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
