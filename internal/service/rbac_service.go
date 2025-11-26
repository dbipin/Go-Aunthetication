package service

import (
	"apiserver/internal/models"
	"apiserver/internal/repository"
	"errors"

	"github.com/go-playground/validator/v10"
)

type RBACService struct {
	rbacRepo       repository.RBACRepository
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	validate       *validator.Validate
}

// NewRBACService creates a new RBAC service
func NewRBACService(
	rbacRepo repository.RBACRepository,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
) *RBACService {
	return &RBACService{
		rbacRepo:       rbacRepo,
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		validate:       validator.New(),
	}
}

// AssignRoleToUser assigns a role to a user
func (s *RBACService) AssignRoleToUser(req *models.AssignRoleRequest) error {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	// Verify user exists
	_, err := s.userRepo.GetByID(req.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify role exists
	_, err = s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Assign role
	return s.rbacRepo.AssignRoleToUser(req.UserID, req.RoleID)
}

// RemoveRoleFromUser removes a role from a user
func (s *RBACService) RemoveRoleFromUser(userID, roleID int) error {
	return s.rbacRepo.RemoveRoleFromUser(userID, roleID)
}

// GetUserRoles retrieves all roles for a user
func (s *RBACService) GetUserRoles(userID int) ([]models.Role, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.rbacRepo.GetUserRoles(userID)
}

// GetUserPermissions retrieves all permissions for a user
func (s *RBACService) GetUserPermissions(userID int) ([]models.Permission, error) {
	// Verify user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.rbacRepo.GetUserPermissions(userID)
}

// GetUserWithRoles retrieves a user with their roles
func (s *RBACService) GetUserWithRoles(userID int) (*models.UserWithRoles, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	roles, err := s.rbacRepo.GetUserRoles(userID)
	if err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""

	return &models.UserWithRoles{
		User:  *user,
		Roles: roles,
	}, nil
}

// GetUserWithPermissions retrieves a user with all their permissions
func (s *RBACService) GetUserWithPermissions(userID int) (*models.UserWithPermissions, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	permissions, err := s.rbacRepo.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	// Clear password
	user.Password = ""

	return &models.UserWithPermissions{
		User:        *user,
		Permissions: permissions,
	}, nil
}

// UserHasRole checks if a user has a specific role
func (s *RBACService) UserHasRole(userID int, roleName string) (bool, error) {
	return s.rbacRepo.UserHasRole(userID, roleName)
}

// UserHasPermission checks if a user has a specific permission
func (s *RBACService) UserHasPermission(userID int, permissionName string) (bool, error) {
	return s.rbacRepo.UserHasPermission(userID, permissionName)
}

// AssignPermissionToRole assigns a permission to a role
func (s *RBACService) AssignPermissionToRole(req *models.AssignPermissionRequest) error {
	// Validate request
	if err := s.validate.Struct(req); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	// Verify role exists
	_, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return errors.New("role not found")
	}

	// Verify permission exists
	_, err = s.permissionRepo.GetByID(req.PermissionID)
	if err != nil {
		return errors.New("permission not found")
	}

	// Assign permission
	return s.rbacRepo.AssignPermissionToRole(req.RoleID, req.PermissionID)
}

// RemovePermissionFromRole removes a permission from a role
func (s *RBACService) RemovePermissionFromRole(roleID, permissionID int) error {
	return s.rbacRepo.RemovePermissionFromRole(roleID, permissionID)
}
