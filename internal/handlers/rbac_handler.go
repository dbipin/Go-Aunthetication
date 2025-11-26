package handlers

import (
	"apiserver/internal/models"
	"apiserver/internal/service"
	"apiserver/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type RBACHandler struct {
	service *service.RBACService
}

// NewRBACHandler creates a new RBAC handler
func NewRBACHandler(service *service.RBACService) *RBACHandler {
	return &RBACHandler{service: service}
}

// AssignRoleToUser handles assigning a role to a user
func (h *RBACHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req models.AssignRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AssignRoleToUser(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Role assigned successfully",
	})
}

// RemoveRoleFromUser handles removing a role from a user
func (h *RBACHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	roleIDStr := chi.URLParam(r, "roleId")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	if err := h.service.RemoveRoleFromUser(userID, roleID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Role removed successfully",
	})
}

// GetUserRoles retrieves all roles for a user
func (h *RBACHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	roles, err := h.service.GetUserRoles(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, roles)
}

// GetUserPermissions retrieves all permissions for a user
func (h *RBACHandler) GetUserPermissions(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	permissions, err := h.service.GetUserPermissions(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, permissions)
}

// GetUserWithRoles retrieves a user with their roles
func (h *RBACHandler) GetUserWithRoles(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userWithRoles, err := h.service.GetUserWithRoles(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, userWithRoles)
}

// GetUserWithPermissions retrieves a user with all their permissions
func (h *RBACHandler) GetUserWithPermissions(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	userWithPerms, err := h.service.GetUserWithPermissions(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, userWithPerms)
}

// AssignPermissionToRole handles assigning a permission to a role
func (h *RBACHandler) AssignPermissionToRole(w http.ResponseWriter, r *http.Request) {
	var req models.AssignPermissionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AssignPermissionToRole(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Permission assigned successfully",
	})
}

// RemovePermissionFromRole handles removing a permission from a role
func (h *RBACHandler) RemovePermissionFromRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := chi.URLParam(r, "roleId")
	permIDStr := chi.URLParam(r, "permissionId")

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	permID, err := strconv.Atoi(permIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	if err := h.service.RemovePermissionFromRole(roleID, permID); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Permission removed successfully",
	})
}
