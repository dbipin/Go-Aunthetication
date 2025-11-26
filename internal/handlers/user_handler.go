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

// Update UserHandler to have roleService

type UserHandler struct {
	service     *service.UserService
	rbacService *service.RBACService
	roleService *service.RoleService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *service.UserService, rbacService *service.RBACService, roleService *service.RoleService) *UserHandler {
	return &UserHandler{
		service:     service,
		rbacService: rbacService,
		roleService: roleService,
	}
}

// Register handles user registration
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// AUTO-ASSIGN "user" ROLE BY NAME (BETTER!)
	userRole, err := h.roleService.GetRoleByName("user")
	if err == nil && userRole != nil {
		_ = h.rbacService.AssignRoleToUser(&models.AssignRoleRequest{
			UserID: user.ID,
			RoleID: userRole.ID,
		})
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Get user roles
	roles, err := h.rbacService.GetUserRoles(user.ID)
	if err != nil {
		roles = []models.Role{}
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
		Roles: roles,
	}

	utils.SuccessResponse(w, http.StatusCreated, response)
}

// Login handles user login
// internal/handlers/user_handler.go

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.Login(&req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// GET USER ROLES (ADD THIS)
	roles, err := h.rbacService.GetUserRoles(user.ID)
	if err != nil {
		roles = []models.Role{} // Empty if error
	}

	response := models.LoginResponse{
		Token: token,
		User:  *user,
		Roles: roles, // ADD THIS
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

// GetMe returns current authenticated user
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(int)

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, user)
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, user)
}

// UpdateUser updates current authenticated user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.UpdateUser(userID, &req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(w, http.StatusOK, user)
}

// DeleteUser deletes current authenticated user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	if err := h.service.DeleteUser(userID); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "User deleted successfully",
	})
}

// GetAllUsers retrieves all users (admin endpoint)
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, users)
}
