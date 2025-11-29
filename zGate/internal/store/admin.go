package store

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Admin represents an admin user in the system
type Admin struct {
	Username     string
	PasswordHash []byte
	Name         string
	Email        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoginAt  *time.Time
}

// CreateAdmin creates a new admin user with a hashed password
func (s *Store) CreateAdmin(username, password, name, email string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO admins (username, password_hash, name, email, created_at, updated_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`

	_, err = s.db.Exec(query, username, hashedPassword, name, email)
	if err != nil {
		return fmt.Errorf("failed to create admin: %w", err)
	}

	return nil
}

// GetAdmin retrieves an admin by username
func (s *Store) GetAdmin(username string) (*Admin, error) {
	query := `
		SELECT username, password_hash, name, email, created_at, updated_at, last_login_at
		FROM admins
		WHERE username = ?
	`

	admin := &Admin{}
	var lastLoginAt sql.NullTime

	err := s.db.QueryRow(query, username).Scan(
		&admin.Username,
		&admin.PasswordHash,
		&admin.Name,
		&admin.Email,
		&admin.CreatedAt,
		&admin.UpdatedAt,
		&lastLoginAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("admin not found: %s", username)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get admin: %w", err)
	}

	if lastLoginAt.Valid {
		admin.LastLoginAt = &lastLoginAt.Time
	}

	return admin, nil
}

// ListAdmins retrieves all admin users
func (s *Store) ListAdmins() ([]*Admin, error) {
	query := `
		SELECT username, password_hash, name, email, created_at, updated_at, last_login_at
		FROM admins
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to list admins: %w", err)
	}
	defer rows.Close()

	var admins []*Admin
	for rows.Next() {
		admin := &Admin{}
		var lastLoginAt sql.NullTime

		err := rows.Scan(
			&admin.Username,
			&admin.PasswordHash,
			&admin.Name,
			&admin.Email,
			&admin.CreatedAt,
			&admin.UpdatedAt,
			&lastLoginAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan admin: %w", err)
		}

		if lastLoginAt.Valid {
			admin.LastLoginAt = &lastLoginAt.Time
		}

		admins = append(admins, admin)
	}

	return admins, nil
}

// DeleteAdmin removes an admin user by username
func (s *Store) DeleteAdmin(username string) error {
	query := `DELETE FROM admins WHERE username = ?`

	result, err := s.db.Exec(query, username)
	if err != nil {
		return fmt.Errorf("failed to delete admin: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("admin not found: %s", username)
	}

	return nil
}

// VerifyAdminPassword checks if the provided password matches the admin's hashed password
func (s *Store) VerifyAdminPassword(username, password string) (bool, error) {
	admin, err := s.GetAdmin(username)
	if err != nil {
		return false, err
	}

	err = bcrypt.CompareHashAndPassword(admin.PasswordHash, []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify password: %w", err)
	}

	return true, nil
}

// UpdateAdminLastLogin updates the last login timestamp for an admin
func (s *Store) UpdateAdminLastLogin(username string) error {
	query := `UPDATE admins SET last_login_at = CURRENT_TIMESTAMP WHERE username = ?`

	_, err := s.db.Exec(query, username)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}
