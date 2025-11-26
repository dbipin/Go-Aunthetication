package repository

import (
	"apiserver/internal/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

// RoleRepository interface defines database operations for roles
type RoleRepository interface {
	Create(role *models.Role) error
	GetByID(id int) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	GetAll() ([]models.Role, error)
	Update(role *models.Role) error
	Delete(id int) error
	GetRolePermissions(roleID int) ([]models.Permission, error)
	GetRoleUsers(roleID int) ([]models.User, error)
}

type roleRepository struct {
	db *sqlx.DB
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create inserts a new role into the database
func (r *roleRepository) Create(role *models.Role) error {
	query := `
		INSERT INTO roles (role_name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(query, role.RoleName, role.Description).
		Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a role by ID
func (r *roleRepository) GetByID(id int) (*models.Role, error) {
	var role models.Role
	query := `
		SELECT id, role_name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	err := r.db.Get(&role, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

// GetByName retrieves a role by name
func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	query := `
		SELECT id, role_name, description, created_at, updated_at
		FROM roles
		WHERE role_name = $1
	`

	err := r.db.Get(&role, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

// GetAll retrieves all roles
func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	query := `
		SELECT id, role_name, description, created_at, updated_at
		FROM roles
		ORDER BY role_name ASC
	`

	err := r.db.Select(&roles, query)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// Update modifies an existing role
func (r *roleRepository) Update(role *models.Role) error {
	query := `
		UPDATE roles
		SET role_name = $1, description = $2
		WHERE id = $3
		RETURNING updated_at
	`

	err := r.db.QueryRow(query, role.RoleName, role.Description, role.ID).
		Scan(&role.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("role not found")
		}
		return err
	}

	return nil
}

// Delete removes a role from the database
func (r *roleRepository) Delete(id int) error {
	query := `DELETE FROM roles WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("role not found")
	}

	return nil
}

// GetRolePermissions retrieves all permissions assigned to a role
func (r *roleRepository) GetRolePermissions(roleID int) ([]models.Permission, error) {
	var permissions []models.Permission
	query := `
		SELECT p.id, p.permission_name, p.resource, p.action, p.description, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action
	`

	err := r.db.Select(&permissions, query, roleID)
	if err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetRoleUsers retrieves all users with this role
func (r *roleRepository) GetRoleUsers(roleID int) ([]models.User, error) {
	var users []models.User
	query := `
		SELECT u.id, u.email, u.name, u.created_at, u.updated_at
		FROM users u
		INNER JOIN user_roles ur ON u.id = ur.user_id
		WHERE ur.role_id = $1
		ORDER BY u.name
	`

	err := r.db.Select(&users, query, roleID)
	if err != nil {
		return nil, err
	}

	// Clear passwords
	for i := range users {
		users[i].Password = ""
	}

	return users, nil
}
