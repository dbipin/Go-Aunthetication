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

	// Initialize layers (Dependency Injection)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

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
		// Public routes (no authentication required)
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			// User routes
			r.Get("/me", userHandler.GetMe)
			r.Put("/me", userHandler.UpdateUser)
			r.Delete("/me", userHandler.DeleteUser)

			// Admin routes (get all users, get user by ID)
			r.Get("/users", userHandler.GetAllUsers)
			r.Get("/users/{id}", userHandler.GetUser)
		})
	})

	// Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf("🚀 Server starting on http://localhost%s", serverAddr)
	log.Printf("📚 API Documentation:")
	log.Printf("   POST   /api/v1/register      - Register new user")
	log.Printf("   POST   /api/v1/login         - Login user")
	log.Printf("   GET    /api/v1/me            - Get current user (auth required)")
	log.Printf("   PUT    /api/v1/me            - Update current user (auth required)")
	log.Printf("   DELETE /api/v1/me            - Delete current user (auth required)")
	log.Printf("   GET    /api/v1/users         - Get all users (auth required)")
	log.Printf("   GET    /api/v1/users/{id}    - Get user by ID (auth required)")

	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
