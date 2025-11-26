package main

import (
	"apiserver/internal/config"
	"apiserver/internal/handlers"
	"apiserver/internal/middleware"
	"apiserver/internal/repository"
	"apiserver/internal/service"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✅ Connected to database successfully!")

	// Initialize repositories (Dependency Injection)
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	rbacRepo := repository.NewRBACRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	roleService := service.NewRoleService(roleRepo)
	permissionService := service.NewPermissionService(permissionRepo)
	rbacService := service.NewRBACService(rbacRepo, userRepo, roleRepo, permissionRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService, rbacService, roleService)

	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)
	rbacHandler := handlers.NewRBACHandler(rbacService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"Server is running"}`))
	})

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// =====================================================================
		// PUBLIC ROUTES (No authentication required)
		// =====================================================================
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		// =====================================================================
		// AUTHENTICATED USER ROUTES (Any logged-in user)
		// =====================================================================
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			// User's own profile
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateUser)
			r.Delete("/me", userHandler.DeleteUser)
		})

		// =====================================================================
		// ADMIN ONLY ROUTES (Requires admin role)
		// =====================================================================
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Use(middleware.RequireRole(rbacRepo, "admin"))

			// -----------------------------------------------------------------
			// USER MANAGEMENT (Admin only)
			// -----------------------------------------------------------------
			r.Get("/users", userHandler.GetAllUsers)
			r.Get("/users/{id}", userHandler.GetUser)

			// -----------------------------------------------------------------
			// ROLE MANAGEMENT (Admin only)
			// -----------------------------------------------------------------
			r.Route("/roles", func(r chi.Router) {
				r.Get("/", roleHandler.GetAllRoles)
				r.Post("/", roleHandler.CreateRole)
				r.Get("/{id}", roleHandler.GetRole)
				r.Put("/{id}", roleHandler.UpdateRole)
				r.Delete("/{id}", roleHandler.DeleteRole)
				r.Get("/{id}/permissions", roleHandler.GetRoleWithPermissions)
				r.Get("/{id}/users", roleHandler.GetRoleWithUsers)
			})

			// -----------------------------------------------------------------
			// PERMISSION MANAGEMENT (Admin only)
			// -----------------------------------------------------------------
			r.Route("/permissions", func(r chi.Router) {
				r.Get("/", permissionHandler.GetAllPermissions)
				r.Post("/", permissionHandler.CreatePermission)
				r.Get("/{id}", permissionHandler.GetPermission)
				r.Put("/{id}", permissionHandler.UpdatePermission)
				r.Delete("/{id}", permissionHandler.DeletePermission)
				r.Get("/{id}/roles", permissionHandler.GetPermissionRoles)
			})

			// -----------------------------------------------------------------
			// RBAC ASSIGNMENTS (Admin only)
			// -----------------------------------------------------------------
			r.Post("/assign-role", rbacHandler.AssignRoleToUser)
			r.Delete("/users/{userId}/roles/{roleId}", rbacHandler.RemoveRoleFromUser)
			r.Get("/users/{userId}/roles", rbacHandler.GetUserRoles)
			r.Get("/users/{userId}/permissions", rbacHandler.GetUserPermissions)
			r.Get("/users/{userId}/with-roles", rbacHandler.GetUserWithRoles)
			r.Get("/users/{userId}/with-permissions", rbacHandler.GetUserWithPermissions)
			r.Post("/assign-permission", rbacHandler.AssignPermissionToRole)
			r.Delete("/roles/{roleId}/permissions/{permissionId}", rbacHandler.RemovePermissionFromRole)
		})
	})

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("🚀 Server starting on http://localhost%s", serverAddr)
	log.Printf("📚 API Documentation:")
	log.Println()
	log.Println("=================================================================")
	log.Println("PUBLIC ROUTES:")
	log.Println("=================================================================")
	log.Printf("   POST   /api/v1/register                 - Register new user")
	log.Printf("   POST   /api/v1/login                    - Login user")
	log.Println()
	log.Println("=================================================================")
	log.Println("USER ROUTES (Auth Required):")
	log.Println("=================================================================")
	log.Printf("   GET    /api/v1/me                       - Get current user")
	log.Printf("   PUT    /api/v1/me                       - Update current user")
	log.Printf("   DELETE /api/v1/me                       - Delete current user")
	log.Println()
	log.Println("=================================================================")
	log.Println("ADMIN ONLY ROUTES (Requires 'admin' role):")
	log.Println("=================================================================")
	log.Println()
	log.Println("USER MANAGEMENT:")
	log.Printf("   GET    /api/v1/users                    - Get all users")
	log.Printf("   GET    /api/v1/users/{id}               - Get user by ID")
	log.Println()
	log.Println("ROLE MANAGEMENT:")
	log.Printf("   GET    /api/v1/roles                    - Get all roles")
	log.Printf("   POST   /api/v1/roles                    - Create role")
	log.Printf("   GET    /api/v1/roles/{id}               - Get role by ID")
	log.Printf("   PUT    /api/v1/roles/{id}               - Update role")
	log.Printf("   DELETE /api/v1/roles/{id}               - Delete role")
	log.Printf("   GET    /api/v1/roles/{id}/permissions   - Get role with permissions")
	log.Printf("   GET    /api/v1/roles/{id}/users         - Get role with users")
	log.Println()
	log.Println("PERMISSION MANAGEMENT:")
	log.Printf("   GET    /api/v1/permissions              - Get all permissions")
	log.Printf("   POST   /api/v1/permissions              - Create permission")
	log.Printf("   GET    /api/v1/permissions/{id}         - Get permission by ID")
	log.Printf("   PUT    /api/v1/permissions/{id}         - Update permission")
	log.Printf("   DELETE /api/v1/permissions/{id}         - Delete permission")
	log.Printf("   GET    /api/v1/permissions/{id}/roles   - Get roles with permission")
	log.Println()
	log.Println("RBAC ASSIGNMENTS:")
	log.Printf("   POST   /api/v1/assign-role                                  - Assign role to user")
	log.Printf("   DELETE /api/v1/users/{userId}/roles/{roleId}                - Remove role from user")
	log.Printf("   GET    /api/v1/users/{userId}/roles                         - Get user's roles")
	log.Printf("   GET    /api/v1/users/{userId}/permissions                   - Get user's permissions")
	log.Printf("   GET    /api/v1/users/{userId}/with-roles                    - Get user with roles")
	log.Printf("   GET    /api/v1/users/{userId}/with-permissions              - Get user with permissions")
	log.Printf("   POST   /api/v1/assign-permission                            - Assign permission to role")
	log.Printf("   DELETE /api/v1/roles/{roleId}/permissions/{permissionId}    - Remove permission from role")
	log.Println()
	log.Println("=================================================================")
	log.Println("🔐 Default Admin Credentials:")
	log.Println("   Email: admin@agmail.com")
	log.Println("   Password: admin123")
	log.Println("=================================================================")
	log.Println()

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
