package models

import "time"

// Role represents a user role in the system
type Role struct {
	ID          int       `json:"id" db:"id"`
	RoleName    string    `json:"role_name" db:"role_name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Permission represents a system permission
type Permission struct {
	ID             int       `json:"id" db:"id"`
	PermissionName string    `json:"permission_name" db:"permission_name"`
	Resource       string    `json:"resource" db:"resource"`
	Action         string    `json:"action" db:"action"`
	Description    string    `json:"description" db:"description"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// UserRole represents the junction table between users and roles
type UserRole struct {
	UserID     int       `json:"user_id" db:"user_id"`
	RoleID     int       `json:"role_id" db:"role_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
}

// RolePermission represents the junction table between roles and permissions
type RolePermission struct {
	RoleID       int       `json:"role_id" db:"role_id"`
	PermissionID int       `json:"permission_id" db:"permission_id"`
	AssignedAt   time.Time `json:"assigned_at" db:"assigned_at"`
}

// Request/Response DTOs

// CreateRoleRequest - Request body for creating a role
type CreateRoleRequest struct {
	RoleName    string `json:"role_name" validate:"required,min=2,max=100"`
	Description string `json:"description" validate:"omitempty,max=500"`
}

// UpdateRoleRequest - Request body for updating a role
type UpdateRoleRequest struct {
	RoleName    string `json:"role_name,omitempty" validate:"omitempty,min=2,max=100"`
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// CreatePermissionRequest - Request body for creating a permission
type CreatePermissionRequest struct {
	PermissionName string `json:"permission_name" validate:"required,min=2,max=100"`
	Resource       string `json:"resource" validate:"required,min=2,max=100"`
	Action         string `json:"action" validate:"required,min=2,max=50"`
	Description    string `json:"description" validate:"omitempty,max=500"`
}

// UpdatePermissionRequest - Request body for updating a permission
type UpdatePermissionRequest struct {
	PermissionName string `json:"permission_name,omitempty" validate:"omitempty,min=2,max=100"`
	Resource       string `json:"resource,omitempty" validate:"omitempty,min=2,max=100"`
	Action         string `json:"action,omitempty" validate:"omitempty,min=2,max=50"`
	Description    string `json:"description,omitempty" validate:"omitempty,max=500"`
}

// AssignRoleRequest - Request body for assigning role to user
type AssignRoleRequest struct {
	UserID int `json:"user_id" validate:"required,gt=0"`
	RoleID int `json:"role_id" validate:"required,gt=0"`
}

// AssignPermissionRequest - Request body for assigning permission to role
type AssignPermissionRequest struct {
	RoleID       int `json:"role_id" validate:"required,gt=0"`
	PermissionID int `json:"permission_id" validate:"required,gt=0"`
}

// Complex response types

// UserWithRoles - User with their assigned roles
type UserWithRoles struct {
	User
	Roles []Role `json:"roles"`
}

// UserWithPermissions - User with all their permissions (through roles)
type UserWithPermissions struct {
	User
	Permissions []Permission `json:"permissions"`
}

// RoleWithPermissions - Role with its assigned permissions
type RoleWithPermissions struct {
	Role
	Permissions []Permission `json:"permissions"`
}

// RoleWithUsers - Role with users who have this role
type RoleWithUsers struct {
	Role
	Users []User `json:"users"`
}
