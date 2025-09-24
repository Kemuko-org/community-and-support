// Kemuko Support Service API
//
// @title Kemuko Support Service API
// @version 1.0.0
// @description EdTech support and ticketing system API for Kemuko educational platform
//
// @contact.name Kemuko Support Team
// @contact.email support@kemuko.com
//
// @license.name MIT
//
// @host localhost:8080
// @basePath /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer token. Format: Bearer {token}
package main

import (
	"log"
	"net/http"
	"fmt"

	"community-support-service/internal/config"
	"community-support-service/internal/database"
	"community-support-service/internal/handlers"
	"community-support-service/internal/repositories/postgres"
	"community-support-service/internal/services"
	"community-support-service/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "community-support-service/docs"
)

func main() {
	// Load .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables and defaults")
	}

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

	// Setup repositories and services
	repo := postgres.NewRepository(db)
	notificationService := services.NewNotificationService(cfg)
	handlers.SetDependencies(repo, notificationService)

	// Setup JWT middleware
	jwtMiddleware := middleware.NewJWTMiddleware(cfg.Auth.JWTSecret)

	// Setup router
	router := mux.NewRouter()
	
	// Health check endpoint (no auth required)
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	// Swagger documentation endpoint
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// Public routes (no authentication required)
	public := api.PathPrefix("/public").Subrouter()
	public.HandleFunc("/categories", handlers.GetCategories).Methods("GET")
	
	// Protected routes (authentication required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(jwtMiddleware.ValidateToken(next.ServeHTTP))
	})
	
	// Student routes
	protected.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handlers.GetStudentTickets(w, r)
		} else if r.Method == "POST" {
			handlers.CreateTicket(w, r)
		}
	}).Methods("GET", "POST")
	protected.HandleFunc("/tickets/{id}", handlers.GetTicketByID).Methods("GET")
	protected.HandleFunc("/tickets/{id}/complete", handlers.CompleteTicket).Methods("POST")
	protected.HandleFunc("/tickets/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handlers.GetTicketComments(w, r)
		} else if r.Method == "POST" {
			handlers.AddTicketComment(w, r)
		}
	}).Methods("GET", "POST")
	
	// Instructor/Admin routes
	instructorRoutes := protected.PathPrefix("/instructor").Subrouter()
	instructorRoutes.Use(func(next http.Handler) http.Handler {
		return middleware.RequireRole("instructor", "admin")(next.ServeHTTP)
	})
	
	instructorRoutes.HandleFunc("/tickets", handlers.GetInstructorTickets).Methods("GET")
	instructorRoutes.HandleFunc("/tickets/{id}", handlers.UpdateTicket).Methods("PUT")
	instructorRoutes.HandleFunc("/tickets/{id}/assign", handlers.AssignTicket).Methods("POST")
	
	// Slack integration endpoints
	slackRoutes := protected.PathPrefix("/slack").Subrouter()
	slackRoutes.Use(func(next http.Handler) http.Handler {
		return middleware.RequireRole("admin", "instructor")(next.ServeHTTP)
	})
	
	slackRoutes.HandleFunc("/reply", handlers.SlackReply).Methods("POST")
	
	// Slack webhook (no auth required - Slack will verify)
	api.HandleFunc("/slack/webhook", handlers.SlackWebhook).Methods("POST")

	// Start server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	log.Printf("Environment: %s", cfg.Server.Env)
	
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}