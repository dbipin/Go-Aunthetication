// internal/middleware/permission.go

package middleware

import (
	"apiserver/internal/repository"
	"apiserver/internal/utils"
	"net/http"
)

// RequireRole checks if user has specific role
func RequireRole(rbacRepo repository.RBACRepository, roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(int)

			hasRole, err := rbacRepo.UserHasRole(userID, roleName)
			if err != nil || !hasRole {
				utils.ErrorResponse(w, http.StatusForbidden, "Access denied: requires "+roleName+" role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePermission checks if user has specific permission
func RequirePermission(rbacRepo repository.RBACRepository, permissionName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(int)

			hasPermission, err := rbacRepo.UserHasPermission(userID, permissionName)
			if err != nil || !hasPermission {
				utils.ErrorResponse(w, http.StatusForbidden, "Access denied: requires "+permissionName+" permission")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
