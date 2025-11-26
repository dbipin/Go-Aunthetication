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

type RoleHandler struct {
	service *service.RoleService
}

// NewRoleHandler creates a new role handler
func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// CreateRole handles role creation
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req models.CreateRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, err := h.service.CreateRole(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, role)
}

// GetRole retrieves a role by ID
func (h *RoleHandler) GetRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	role, err := h.service.GetRoleByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Role not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, role)
}

// GetAllRoles retrieves all roles
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetAllRoles()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve roles")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, roles)
}

// UpdateRole updates a role
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	var req models.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	role, err := h.service.UpdateRole(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, role)
}

// DeleteRole deletes a role
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	if err := h.service.DeleteRole(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete role")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Role deleted successfully",
	})
}

// GetRoleWithPermissions retrieves a role with its permissions
func (h *RoleHandler) GetRoleWithPermissions(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	roleWithPerms, err := h.service.GetRoleWithPermissions(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Role not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, roleWithPerms)
}

// GetRoleWithUsers retrieves a role with users who have it
func (h *RoleHandler) GetRoleWithUsers(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid role ID")
		return
	}

	roleWithUsers, err := h.service.GetRoleWithUsers(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Role not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, roleWithUsers)
}
