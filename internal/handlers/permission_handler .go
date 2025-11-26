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

type PermissionHandler struct {
	service *service.PermissionService
}

// NewPermissionHandler creates a new permission handler
func NewPermissionHandler(service *service.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

// CreatePermission handles permission creation
func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePermissionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	permission, err := h.service.CreatePermission(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, permission)
}

// GetPermission retrieves a permission by ID
func (h *PermissionHandler) GetPermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	permission, err := h.service.GetPermissionByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Permission not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, permission)
}

// GetAllPermissions retrieves all permissions
func (h *PermissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.service.GetAllPermissions()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve permissions")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, permissions)
}

// UpdatePermission updates a permission
func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	var req models.UpdatePermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	permission, err := h.service.UpdatePermission(id, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, permission)
}

// DeletePermission deletes a permission
func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	if err := h.service.DeletePermission(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete permission")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Permission deleted successfully",
	})
}

// GetPermissionRoles retrieves all roles that have this permission
func (h *PermissionHandler) GetPermissionRoles(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid permission ID")
		return
	}

	roles, err := h.service.GetPermissionRoles(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Permission not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, roles)
}
