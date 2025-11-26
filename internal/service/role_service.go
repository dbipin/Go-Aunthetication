package service

import (
	"apiserver/internal/models"
	"apiserver/internal/repository"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type RoleService struct {
	repo     repository.RoleRepository
	validate *validator.Validate
}

// NewRoleService creates a new role service
func NewRoleService(repo repository.RoleRepository) *RoleService {
	return &RoleService{
		repo:     repo,
		validate: validator.New(),
	}
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(req *models.CreateRoleRequest) (*models.Role, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Normalize role name
	req.RoleName = strings.ToLower(strings.TrimSpace(req.RoleName))

	// Check if role already exists
	existingRole, _ := s.repo.GetByName(req.RoleName)
	if existingRole != nil {
		return nil, errors.New("role already exists")
	}

	// Create role
	role := &models.Role{
		RoleName:    req.RoleName,
		Description: req.Description,
	}

	if err := s.repo.Create(role); err != nil {
		return nil, errors.New("failed to create role: " + err.Error())
	}

	return role, nil
}

// GetRoleByID retrieves a role by ID
func (s *RoleService) GetRoleByID(id int) (*models.Role, error) {
	return s.repo.GetByID(id)
}

// GetRoleByName retrieves a role by name
func (s *RoleService) GetRoleByName(name string) (*models.Role, error) {
	return s.repo.GetByName(strings.ToLower(name))
}

// GetAllRoles retrieves all roles
func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.repo.GetAll()
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(id int, req *models.UpdateRoleRequest) (*models.Role, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Get existing role
	role, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.RoleName != "" {
		req.RoleName = strings.ToLower(strings.TrimSpace(req.RoleName))

		// Check if new name is already taken by another role
		existingRole, _ := s.repo.GetByName(req.RoleName)
		if existingRole != nil && existingRole.ID != id {
			return nil, errors.New("role name already in use")
		}

		role.RoleName = req.RoleName
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	// Update in database
	if err := s.repo.Update(role); err != nil {
		return nil, errors.New("failed to update role: " + err.Error())
	}

	return role, nil
}

// DeleteRole deletes a role
func (s *RoleService) DeleteRole(id int) error {
	return s.repo.Delete(id)
}

// GetRoleWithPermissions retrieves a role with its permissions
func (s *RoleService) GetRoleWithPermissions(roleID int) (*models.RoleWithPermissions, error) {
	role, err := s.repo.GetByID(roleID)
	if err != nil {
		return nil, err
	}

	permissions, err := s.repo.GetRolePermissions(roleID)
	if err != nil {
		return nil, err
	}

	return &models.RoleWithPermissions{
		Role:        *role,
		Permissions: permissions,
	}, nil
}

// GetRoleWithUsers retrieves a role with users who have this role
func (s *RoleService) GetRoleWithUsers(roleID int) (*models.RoleWithUsers, error) {
	role, err := s.repo.GetByID(roleID)
	if err != nil {
		return nil, err
	}

	users, err := s.repo.GetRoleUsers(roleID)
	if err != nil {
		return nil, err
	}

	return &models.RoleWithUsers{
		Role:  *role,
		Users: users,
	}, nil
}
