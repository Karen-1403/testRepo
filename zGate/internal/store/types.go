package store

import "time"

// Permission mirrors the YAML permission structure.
type Permission struct {
	Database string `json:"database"`
	Level    string `json:"level"`
}

// Database represents a backend database definition.
type Database struct {
	Name                 string    `json:"name"`
	Type                 string    `json:"type"`
	Description          string    `json:"description"`
	BackendAddr          string    `json:"backend_addr"`
	AdminUsername        string    `json:"admin_username"`
	AdminPassword        string    `json:"admin_password"`
	AvailablePermissions []string  `json:"available_permissions"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// Role contains description and permissions.
type Role struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

// User captures the credentials and role bindings.
type User struct {
	Username          string       `json:"username"`
	PasswordHash      string       `json:"-"`
	Roles             []string     `json:"roles"`
	CustomPermissions []Permission `json:"custom_permissions"`
	CreatedAt         time.Time    `json:"created_at"`
}

// RefreshToken represents a stored refresh token for session management.
type RefreshToken struct {
	ID         int64      `json:"id"`
	TokenHash  string     `json:"-"`
	Username   string     `json:"username"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt time.Time  `json:"last_used_at"`
	Revoked    bool       `json:"revoked"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`
	UserAgent  string     `json:"user_agent,omitempty"`
	IPAddress  string     `json:"ip_address,omitempty"`
}
