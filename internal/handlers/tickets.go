package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"community-support-service/internal/models"
	"community-support-service/internal/repositories"
	"community-support-service/internal/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	repo *repositories.Repository
	notificationService *services.NotificationService
)

func SetDependencies(r *repositories.Repository, ns *services.NotificationService) {
	repo = r
	notificationService = ns
}

// GetStudentTickets godoc
// @Summary Get student's tickets
// @Description Retrieve all tickets created by the authenticated student
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param status query string false "Filter by ticket status" Enums(open,in_progress,resolved,closed)
// @Param priority query string false "Filter by ticket priority" Enums(low,medium,high,urgent)
// @Param type query string false "Filter by ticket type" Enums(general,technical,course,assignment,grading,platform,content)
// @Param search query string false "Search in title and description"
// @Param page query int false "Page number for pagination" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Router /tickets [get]
func GetStudentTickets(w http.ResponseWriter, r *http.Request) {
	studentID := r.Header.Get("X-User-ID")
	if studentID == "" {
		http.Error(w, "User ID not found in request", http.StatusUnauthorized)
		return
	}

	filters := repositories.TicketFilters{}
	if status := r.URL.Query().Get("status"); status != "" {
		filters.Status = &status
	}
	if priority := r.URL.Query().Get("priority"); priority != "" {
		filters.Priority = &priority
	}
	if ticketType := r.URL.Query().Get("type"); ticketType != "" {
		filters.Type = &ticketType
	}
	if search := r.URL.Query().Get("search"); search != "" {
		filters.Search = &search
	}

	tickets, err := repo.Ticket.GetByStudentID(r.Context(), studentID, filters)
	if err != nil {
		http.Error(w, "Failed to fetch tickets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tickets": tickets,
		"total":   len(tickets),
	})
}

// CreateTicket godoc
// @Summary Create new ticket
// @Description Create a new support ticket
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /tickets [post]
func CreateTicket(w http.ResponseWriter, r *http.Request) {
	studentID := r.Header.Get("X-User-ID")
	studentEmail := r.Header.Get("X-User-Email")
	if studentID == "" || studentEmail == "" {
		http.Error(w, "User information not found in request", http.StatusUnauthorized)
		return
	}

	var req models.CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ticketNumber := generateTicketNumber()
	now := time.Now()
	
	ticket := &models.Ticket{
		ID:           uuid.New(),
		TicketNumber: ticketNumber,
		Title:        req.Title,
		Description:  req.Description,
		Status:       models.TicketStatusOpen,
		Priority:     req.Priority,
		Type:         req.Type,
		StudentID:    studentID,
		CourseID:     req.CourseID,
		CategoryID:   req.CategoryID,
		Metadata:     req.Metadata,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := repo.Ticket.Create(r.Context(), ticket); err != nil {
		http.Error(w, "Failed to create ticket", http.StatusInternalServerError)
		return
	}

	history := &models.TicketHistory{
		ID:          uuid.New(),
		TicketID:    ticket.ID,
		UserID:      studentID,
		Action:      "created",
		Description: StringPtr("Ticket created by student"),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
	}
	if err := repo.History.Create(r.Context(), history); err != nil {
		fmt.Printf("Failed to create ticket history: %v\n", err)
	}

	if err := notificationService.SendTicketCreatedNotifications(r.Context(), ticket, studentEmail); err != nil {
		fmt.Printf("Failed to send notifications: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ticket": ticket,
		"message": "Ticket created successfully",
	})
}

// GetTicketByID godoc
// @Summary Get ticket by ID
// @Description Retrieve a specific ticket by its ID
// @Tags tickets
// @Security BearerAuth
// @Produce json
// @Param id path string true "Ticket ID" Format(uuid)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/{id} [get]
func GetTicketByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	ticket, err := repo.Ticket.GetByID(r.Context(), ticketID)
	if err != nil {
		http.Error(w, "Failed to fetch ticket", http.StatusInternalServerError)
		return
	}
	if ticket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	studentID := r.Header.Get("X-User-ID")
	userRole := r.Header.Get("X-User-Role")
	if ticket.StudentID != studentID && userRole != "instructor" && userRole != "admin" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ticket": ticket,
	})
}

// CompleteTicket godoc
// @Summary Complete ticket
// @Description Mark a ticket as resolved/completed
// @Tags tickets
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID" Format(uuid)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/{id}/complete [post]
func CompleteTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ticketID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("X-User-ID")
	userEmail := r.Header.Get("X-User-Email")
	userRole := r.Header.Get("X-User-Role")
	if userID == "" || userRole == "" {
		http.Error(w, "User information not found in request", http.StatusUnauthorized)
		return
	}

	ticket, err := repo.Ticket.GetByID(r.Context(), ticketID)
	if err != nil {
		http.Error(w, "Failed to fetch ticket", http.StatusInternalServerError)
		return
	}
	if ticket == nil {
		http.Error(w, "Ticket not found", http.StatusNotFound)
		return
	}

	if userRole != "instructor" && userRole != "admin" {
		http.Error(w, "Only instructors and admins can complete tickets", http.StatusForbidden)
		return
	}

	if ticket.Status == models.TicketStatusResolved || ticket.Status == models.TicketStatusClosed {
		http.Error(w, "Ticket is already completed", http.StatusBadRequest)
		return
	}

	now := time.Now()
	oldStatus := string(ticket.Status)
	ticket.Status = models.TicketStatusResolved
	ticket.UpdatedAt = now
	ticket.ResolvedAt = &now

	if err := repo.Ticket.Update(r.Context(), ticket); err != nil {
		http.Error(w, "Failed to update ticket", http.StatusInternalServerError)
		return
	}

	history := &models.TicketHistory{
		ID:          uuid.New(),
		TicketID:    ticket.ID,
		UserID:      userID,
		Action:      "completed",
		OldValue:    &oldStatus,
		NewValue:    StringPtr(string(models.TicketStatusResolved)),
		Description: StringPtr(fmt.Sprintf("Ticket completed by %s", userRole)),
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
	}
	if err := repo.History.Create(r.Context(), history); err != nil {
		fmt.Printf("Failed to create ticket history: %v\n", err)
	}

	if err := sendTicketCompletedNotifications(r.Context(), ticket, userEmail); err != nil {
		fmt.Printf("Failed to send completion notifications: %v\n", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ticket":  ticket,
		"message": "Ticket completed successfully",
	})
}

func generateTicketNumber() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("KEMUKO-%d", timestamp)
}

func StringPtr(s string) *string {
	return &s
}

func sendTicketCompletedNotifications(ctx context.Context, ticket *models.Ticket, resolverEmail string) error {
	slackReq := &models.SlackNotificationRequest{
		Channel: "#kemuko-support",
		Message: fmt.Sprintf("✅ Ticket Completed"),
		Blocks: []models.SlackBlock{
			{
				Type: "section",
				Text: &models.SlackText{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*✅ Ticket Completed*\n*Ticket:* #%s\n*Title:* %s\n*Completed by:* %s\n*Status:* %s",
						ticket.TicketNumber, ticket.Title, resolverEmail, ticket.Status),
				},
			},
		},
		TemplateData: map[string]interface{}{
			"ticketId":     ticket.ID.String(),
			"ticketNumber": ticket.TicketNumber,
			"status":       string(ticket.Status),
		},
	}

	return notificationService.SendSlackRequest(ctx, slackReq)
}