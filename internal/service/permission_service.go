package service

import (
	"apiserver/internal/models"
	"apiserver/internal/repository"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type PermissionService struct {
	repo     repository.PermissionRepository
	validate *validator.Validate
}

// NewPermissionService creates a new permission service
func NewPermissionService(repo repository.PermissionRepository) *PermissionService {
	return &PermissionService{
		repo:     repo,
		validate: validator.New(),
	}
}

// CreatePermission creates a new permission
func (s *PermissionService) CreatePermission(req *models.CreatePermissionRequest) (*models.Permission, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Normalize fields
	req.PermissionName = strings.ToLower(strings.TrimSpace(req.PermissionName))
	req.Resource = strings.ToLower(strings.TrimSpace(req.Resource))
	req.Action = strings.ToLower(strings.TrimSpace(req.Action))

	// Check if permission already exists
	existingPermission, _ := s.repo.GetByName(req.PermissionName)
	if existingPermission != nil {
		return nil, errors.New("permission already exists")
	}

	// Create permission
	permission := &models.Permission{
		PermissionName: req.PermissionName,
		Resource:       req.Resource,
		Action:         req.Action,
		Description:    req.Description,
	}

	if err := s.repo.Create(permission); err != nil {
		return nil, errors.New("failed to create permission: " + err.Error())
	}

	return permission, nil
}

// GetPermissionByID retrieves a permission by ID
func (s *PermissionService) GetPermissionByID(id int) (*models.Permission, error) {
	return s.repo.GetByID(id)
}

// GetPermissionByName retrieves a permission by name
func (s *PermissionService) GetPermissionByName(name string) (*models.Permission, error) {
	return s.repo.GetByName(strings.ToLower(name))
}

// GetAllPermissions retrieves all permissions
func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	return s.repo.GetAll()
}

// UpdatePermission updates an existing permission
func (s *PermissionService) UpdatePermission(id int, req *models.UpdatePermissionRequest) (*models.Permission, error) {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	// Get existing permission
	permission, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.PermissionName != "" {
		req.PermissionName = strings.ToLower(strings.TrimSpace(req.PermissionName))

		// Check if new name is already taken by another permission
		existingPermission, _ := s.repo.GetByName(req.PermissionName)
		if existingPermission != nil && existingPermission.ID != id {
			return nil, errors.New("permission name already in use")
		}

		permission.PermissionName = req.PermissionName
	}

	if req.Resource != "" {
		permission.Resource = strings.ToLower(strings.TrimSpace(req.Resource))
	}

	if req.Action != "" {
		permission.Action = strings.ToLower(strings.TrimSpace(req.Action))
	}

	if req.Description != "" {
		permission.Description = req.Description
	}

	// Update in database
	if err := s.repo.Update(permission); err != nil {
		return nil, errors.New("failed to update permission: " + err.Error())
	}

	return permission, nil
}

// DeletePermission deletes a permission
func (s *PermissionService) DeletePermission(id int) error {
	return s.repo.Delete(id)
}

// GetPermissionRoles retrieves all roles that have this permission
func (s *PermissionService) GetPermissionRoles(permissionID int) ([]models.Role, error) {
	return s.repo.GetPermissionRoles(permissionID)
}
