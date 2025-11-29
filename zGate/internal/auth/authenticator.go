package auth

import (
	"fmt"

	"github.com/zGate-Team/zGate-Platform/internal/store"
)

// Authenticator handles user authentication
type Authenticator struct {
	store *store.Store
}

// NewAuthenticator creates a new authenticator
func NewAuthenticator(store *store.Store) *Authenticator {
	return &Authenticator{store: store}
}

// Authenticate validates username and password, returns user with resolved permissions
func (a *Authenticator) Authenticate(username, password string) (*UserWithPermissions, error) {
	if err := a.store.VerifyPassword(username, password); err != nil {
		return nil, fmt.Errorf("authentication failed: invalid credentials")
	}

	user, err := a.store.GetUser(username)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	rolePerms, err := a.store.GetPermissionsForRoles(user.Roles)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve role permissions: %w", err)
	}

	permissions := append(rolePerms, user.CustomPermissions...)

	return &UserWithPermissions{
		Username:    user.Username,
		Roles:       user.Roles,
		Permissions: permissions,
	}, nil
}

// UserWithPermissions represents a user with resolved permissions
type UserWithPermissions struct {
	Username    string
	Roles       []string
	Permissions []store.Permission
}
