package main

import (
	"log"
	"net/http"
	"fmt"

	"community-support-service/internal/config"
	"community-support-service/internal/database"
	"community-support-service/pkg/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Setup JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.Auth.JWTSecret)

	// Setup router
	router := mux.NewRouter()
	
	// Health check endpoint (no auth required)
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"edtech-support-service"}`))
	}).Methods("GET")

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// Public routes (no authentication required)
	public := api.PathPrefix("/public").Subrouter()
	public.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Public categories endpoint - Coming soon"}`))
	}).Methods("GET")
	
	// Protected routes (authentication required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(jwtMiddleware.ValidateToken(next.ServeHTTP))
	})
	
	// Student routes
	protected.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Student tickets endpoint - Coming soon"}`))
	}).Methods("GET", "POST")
	
	// Instructor/Admin routes
	instructorRoutes := protected.PathPrefix("/instructor").Subrouter()
	instructorRoutes.Use(func(next http.Handler) http.Handler {
		return middleware.RequireRole("instructor", "admin")(next.ServeHTTP)
	})
	
	instructorRoutes.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Instructor tickets endpoint - Coming soon"}`))
	}).Methods("GET")

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Environment: %s", cfg.Server.Env)
	
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}