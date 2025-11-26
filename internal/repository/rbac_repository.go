package repository

import (
	"apiserver/internal/models"
	"errors"

	"github.com/jmoiron/sqlx"
)

// RBACRepository interface defines operations for role-based access control
type RBACRepository interface {
	// User-Role operations
	AssignRoleToUser(userID, roleID int) error
	RemoveRoleFromUser(userID, roleID int) error
	GetUserRoles(userID int) ([]models.Role, error)
	GetUserPermissions(userID int) ([]models.Permission, error)
	UserHasRole(userID int, roleName string) (bool, error)
	UserHasPermission(userID int, permissionName string) (bool, error)

	// Role-Permission operations
	AssignPermissionToRole(roleID, permissionID int) error
	RemovePermissionFromRole(roleID, permissionID int) error
}

type rbacRepository struct {
	db *sqlx.DB
}

// NewRBACRepository creates a new RBAC repository
func NewRBACRepository(db *sqlx.DB) RBACRepository {
	return &rbacRepository{db: db}
}

// AssignRoleToUser assigns a role to a user
func (r *rbacRepository) AssignRoleToUser(userID, roleID int) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	_, err := r.db.Exec(query, userID, roleID)
	return err
}

// RemoveRoleFromUser removes a role from a user
func (r *rbacRepository) RemoveRoleFromUser(userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	result, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user-role assignment not found")
	}

	return nil
}

// GetUserRoles retrieves all roles assigned to a user
func (r *rbacRepository) GetUserRoles(userID int) ([]models.Role, error) {
	var roles []models.Role
	query := `
		SELECT r.id, r.role_name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.role_name
	`

	err := r.db.Select(&roles, query, userID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// GetUserPermissions retrieves all permissions a user has (through their roles)
func (r *rbacRepository) GetUserPermissions(userID int) ([]models.Permission, error) {
	var permissions []models.Permission
	query := `
		SELECT DISTINCT p.id, p.permission_name, p.resource, p.action, p.description, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY p.resource, p.action
	`

	err := r.db.Select(&permissions, query, userID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// UserHasRole checks if a user has a specific role
func (r *rbacRepository) UserHasRole(userID int, roleName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_roles ur
			INNER JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 AND r.role_name = $2
		)
	`

	err := r.db.Get(&exists, query, userID, roleName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UserHasPermission checks if a user has a specific permission (through any of their roles)
func (r *rbacRepository) UserHasPermission(userID int, permissionName string) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(
			SELECT 1
			FROM user_roles ur
			INNER JOIN role_permissions rp ON ur.role_id = rp.role_id
			INNER JOIN permissions p ON rp.permission_id = p.id
			WHERE ur.user_id = $1 AND p.permission_name = $2
		)
	`

	err := r.db.Get(&exists, query, userID, permissionName)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// AssignPermissionToRole assigns a permission to a role
func (r *rbacRepository) AssignPermissionToRole(roleID, permissionID int) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES ($1, $2)
		ON CONFLICT (role_id, permission_id) DO NOTHING
	`

	_, err := r.db.Exec(query, roleID, permissionID)
	return err
}

// RemovePermissionFromRole removes a permission from a role
func (r *rbacRepository) RemovePermissionFromRole(roleID, permissionID int) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	result, err := r.db.Exec(query, roleID, permissionID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("role-permission assignment not found")
	}

	return nil
}
