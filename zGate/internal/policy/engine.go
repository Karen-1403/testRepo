package policy

import (
	"github.com/zGate-Team/zGate-Platform/internal/auth"
	"github.com/zGate-Team/zGate-Platform/internal/store"
	"github.com/zGate-Team/zGate-Platform/internal/utils"
)

// Engine is the policy evaluation engine
type Engine struct {
	store *store.Store
}

// NewEngine creates a new policy engine with access to all configs
func NewEngine(store *store.Store) *Engine {
	return &Engine{store: store}
}

// getFreshPermissions looks up the user and calculates permissions from the current config
func (e *Engine) getFreshPermissions(username string) ([]store.Permission, error) {
	user, err := e.store.GetUser(username)
	if err != nil {
		return nil, err
	}

	rolePerms, err := e.store.GetPermissionsForRoles(user.Roles)
	if err != nil {
		return nil, err
	}

	return append(rolePerms, user.CustomPermissions...), nil
}

// CanAccess checks if user can access a database using fresh permissions
func (e *Engine) CanAccess(claims *auth.Claims, databaseName string) bool {
	// Ignore claims.Permissions, lookup fresh permissions
	perms, err := e.getFreshPermissions(claims.Username)
	if err != nil {
		utils.Logger.Warn("user not found during policy check", "username", claims.Username)
		return false
	}

	for _, perm := range perms {
		if perm.Database == databaseName {
			return true
		}
	}
	return false
}

// GetAllowedDatabases returns list of databases user can access using fresh permissions
func (e *Engine) GetAllowedDatabases(claims *auth.Claims) []DatabaseInfo {
	perms, err := e.getFreshPermissions(claims.Username)
	if err != nil {
		utils.Logger.Warn("user not found during list check", "username", claims.Username)
		return []DatabaseInfo{}
	}

	databases, err := e.store.ListDatabases()
	if err != nil {
		utils.Logger.Error("failed to list databases", "error", err)
		return []DatabaseInfo{}
	}

	var allowed []DatabaseInfo
	for _, db := range databases {
		var userPerm *store.Permission
		for i := range perms {
			if perms[i].Database == db.Name {
				userPerm = &perms[i]
				break
			}
		}

		if userPerm != nil {
			allowed = append(allowed, DatabaseInfo{
				Name:        db.Name,
				Type:        db.Type,
				Permissions: userPerm.Level,
				Status:      "online",
				Description: db.Description,
			})
		}
	}

	return allowed
}

// GetPermissionLevel returns the permission level for a database using fresh permissions
func (e *Engine) GetPermissionLevel(claims *auth.Claims, databaseName string) string {
	perms, err := e.getFreshPermissions(claims.Username)
	if err != nil {
		return ""
	}

	for _, perm := range perms {
		if perm.Database == databaseName {
			return perm.Level
		}
	}
	return ""
}

// DatabaseInfo represents database information for CLI display
type DatabaseInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Permissions string `json:"permissions"`
	Status      string `json:"status"`
	Description string `json:"description"`
}
