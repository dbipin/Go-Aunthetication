package repository

import (
	"apiserver/internal/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// PermissionRepository interface defines database operations for permissions
type PermissionRepository interface {
	Create(permission *models.Permission) error
	GetByID(id int) (*models.Permission, error)
	GetByName(name string) (*models.Permission, error)
	GetAll() ([]models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id int) error
	GetPermissionRoles(permissionID int) ([]models.Role, error)
}

type permissionRepository struct {
	db *sqlx.DB
}

// NewPermissionRepository creates a new permission repository
func NewPermissionRepository(db *sqlx.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// Create inserts a new permission into the database
func (r *permissionRepository) Create(permission *models.Permission) error {
	query := `
		INSERT INTO permissions (permission_name, resource, action, description)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, permission.PermissionName, permission.Resource,
		permission.Action, permission.Description).
		Scan(&permission.ID, &permission.CreatedAt, &permission.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a permission by ID
func (r *permissionRepository) GetByID(id int) (*models.Permission, error) {
	var permission models.Permission
	query := `
		SELECT id, permission_name, resource, action, description, created_at, updated_at
		FROM permissions
		WHERE id = $1
	`

	err := r.db.Get(&permission, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}

// GetByName retrieves a permission by name
func (r *permissionRepository) GetByName(name string) (*models.Permission, error) {
	var permission models.Permission
	query := `
		SELECT id, permission_name, resource, action, description, created_at, updated_at
		FROM permissions
		WHERE permission_name = $1
	`

	err := r.db.Get(&permission, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return &permission, nil
}

// GetAll retrieves all permissions
func (r *permissionRepository) GetAll() ([]models.Permission, error) {
	var permissions []models.Permission
	query := `
		SELECT id, permission_name, resource, action, description, created_at, updated_at
		FROM permissions
		ORDER BY resource, action
	`

	err := r.db.Select(&permissions, query)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// Update modifies an existing permission
func (r *permissionRepository) Update(permission *models.Permission) error {
	query := `
		UPDATE permissions
		SET permission_name = $1, resource = $2, action = $3, description = $4
		WHERE id = $5
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, permission.PermissionName, permission.Resource,
		permission.Action, permission.Description, permission.ID).
		Scan(&permission.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("permission not found")
		}
		return err
	}

	return nil
}

// Delete removes a permission from the database
func (r *permissionRepository) Delete(id int) error {
	query := `DELETE FROM permissions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("permission not found")
	}

	return nil
}

// GetPermissionRoles retrieves all roles that have this permission
func (r *permissionRepository) GetPermissionRoles(permissionID int) ([]models.Role, error) {
	var roles []models.Role
	query := `
		SELECT r.id, r.role_name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN role_permissions rp ON r.id = rp.role_id
		WHERE rp.permission_id = $1
		ORDER BY r.role_name
	`

	err := r.db.Select(&roles, query, permissionID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
