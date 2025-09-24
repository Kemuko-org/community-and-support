package handlers

import (
	"net/http"
)

// GetTicketComments godoc
// @Summary Get ticket comments
// @Description Retrieve all comments for a specific ticket
// @Tags comments
// @Security BearerAuth
// @Produce json
// @Param id path string true "Ticket ID" Format(uuid)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/{id}/comments [get]
func GetTicketComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Get ticket comments endpoint - Coming soon"}`))
}

// AddTicketComment godoc
// @Summary Add comment to ticket
// @Description Add a new comment to a specific ticket
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Ticket ID" Format(uuid)
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tickets/{id}/comments [post]
func AddTicketComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message":"Add ticket comment endpoint - Coming soon"}`))
}