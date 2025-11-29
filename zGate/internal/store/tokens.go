package store

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// CreateRefreshToken generates and stores a new refresh token for a user.
func (s *Store) CreateRefreshToken(username, userAgent, ipAddress string, duration time.Duration) (string, error) {
	// Generate cryptographically secure random token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// Hash the token for storage
	tokenHash := hashToken(token)

	// Calculate expiration
	expiresAt := time.Now().Add(duration)

	// Insert into database
	_, err := s.db.Exec(`
		INSERT INTO refresh_tokens (token_hash, username, expires_at, user_agent, ip_address)
		VALUES (?, ?, ?, ?, ?)
	`, tokenHash, username, expiresAt, userAgent, ipAddress)

	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %w", err)
	}

	return token, nil
}

// ValidateRefreshToken checks if a refresh token is valid and returns its details.
func (s *Store) ValidateRefreshToken(token string) (*RefreshToken, error) {
	tokenHash := hashToken(token)

	var rt RefreshToken
	var revokedAt *time.Time

	err := s.db.QueryRow(`
		SELECT id, token_hash, username, expires_at, created_at, last_used_at, revoked, revoked_at, user_agent, ip_address
		FROM refresh_tokens
		WHERE token_hash = ?
	`, tokenHash).Scan(
		&rt.ID,
		&rt.TokenHash,
		&rt.Username,
		&rt.ExpiresAt,
		&rt.CreatedAt,
		&rt.LastUsedAt,
		&rt.Revoked,
		&revokedAt,
		&rt.UserAgent,
		&rt.IPAddress,
	)

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	rt.RevokedAt = revokedAt

	// Check if token is revoked
	if rt.Revoked {
		return nil, fmt.Errorf("refresh token has been revoked")
	}

	// Check if token is expired
	if time.Now().After(rt.ExpiresAt) {
		return nil, fmt.Errorf("refresh token has expired")
	}

	// Update last_used_at timestamp
	_, err = s.db.Exec(`
		UPDATE refresh_tokens
		SET last_used_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, rt.ID)

	if err != nil {
		// Log but don't fail - this is not critical
		// Could use logger here if available
	}

	return &rt, nil
}

// RevokeRefreshToken marks a specific refresh token as revoked.
func (s *Store) RevokeRefreshToken(token string) error {
	tokenHash := hashToken(token)

	result, err := s.db.Exec(`
		UPDATE refresh_tokens
		SET revoked = 1, revoked_at = CURRENT_TIMESTAMP
		WHERE token_hash = ? AND revoked = 0
	`, tokenHash)

	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check revocation: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found or already revoked")
	}

	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a specific user.
func (s *Store) RevokeAllUserTokens(username string) error {
	_, err := s.db.Exec(`
		UPDATE refresh_tokens
		SET revoked = 1, revoked_at = CURRENT_TIMESTAMP
		WHERE username = ? AND revoked = 0
	`, username)

	if err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}

	return nil
}

// CleanupExpiredTokens removes expired and revoked tokens from the database.
func (s *Store) CleanupExpiredTokens() error {
	// Delete tokens that expired more than 24 hours ago
	result, err := s.db.Exec(`
		DELETE FROM refresh_tokens
		WHERE expires_at < datetime('now', '-1 day')
	`)

	if err != nil {
		return fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}

	// Also delete revoked tokens older than 7 days
	_, err = s.db.Exec(`
		DELETE FROM refresh_tokens
		WHERE revoked = 1 AND revoked_at < datetime('now', '-7 day')
	`)

	if err != nil {
		return fmt.Errorf("failed to cleanup revoked tokens: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		// Could log this information
		_ = rowsAffected
	}

	return nil
}

// GetUserActiveLogins returns all active (non-revoked, non-expired) sessions for a user.
func (s *Store) GetUserActiveLogins(username string) ([]RefreshToken, error) {
	rows, err := s.db.Query(`
		SELECT id, token_hash, username, expires_at, created_at, last_used_at, revoked, revoked_at, user_agent, ip_address
		FROM refresh_tokens
		WHERE username = ? AND revoked = 0 AND expires_at > CURRENT_TIMESTAMP
		ORDER BY last_used_at DESC
	`, username)

	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	defer rows.Close()

	var sessions []RefreshToken
	for rows.Next() {
		var rt RefreshToken
		var revokedAt *time.Time

		err := rows.Scan(
			&rt.ID,
			&rt.TokenHash,
			&rt.Username,
			&rt.ExpiresAt,
			&rt.CreatedAt,
			&rt.LastUsedAt,
			&rt.Revoked,
			&revokedAt,
			&rt.UserAgent,
			&rt.IPAddress,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		rt.RevokedAt = revokedAt
		sessions = append(sessions, rt)
	}

	return sessions, nil
}

// RevokeActiveLoginByID revokes a specific session by its ID (for session management UI).
func (s *Store) RevokeActiveLoginByID(sessionID int64, username string) error {
	result, err := s.db.Exec(`
		UPDATE refresh_tokens
		SET revoked = 1, revoked_at = CURRENT_TIMESTAMP
		WHERE id = ? AND username = ? AND revoked = 0
	`, sessionID, username)

	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check revocation: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found or already revoked")
	}

	return nil
}

// GetAllActiveLogins returns all active (non-revoked, non-expired) sessions for all users.
// Intended for administrative visibility; caller must enforce authorization.
func (s *Store) GetAllActiveLogins() ([]RefreshToken, error) {
	rows, err := s.db.Query(`
		SELECT id, token_hash, username, expires_at, created_at, last_used_at, revoked, revoked_at, user_agent, ip_address
		FROM refresh_tokens
		WHERE revoked = 0 AND expires_at > CURRENT_TIMESTAMP
		ORDER BY last_used_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get active logins: %w", err)
	}
	defer rows.Close()

	var sessions []RefreshToken
	for rows.Next() {
		var rt RefreshToken
		var revokedAt *time.Time
		if err := rows.Scan(
			&rt.ID,
			&rt.TokenHash,
			&rt.Username,
			&rt.ExpiresAt,
			&rt.CreatedAt,
			&rt.LastUsedAt,
			&rt.Revoked,
			&revokedAt,
			&rt.UserAgent,
			&rt.IPAddress,
		); err != nil {
			return nil, fmt.Errorf("failed to scan active login: %w", err)
		}
		rt.RevokedAt = revokedAt
		sessions = append(sessions, rt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}
	return sessions, nil
}

// hashToken creates a SHA-256 hash of the token for secure storage.
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
