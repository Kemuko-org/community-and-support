package handlers

import (
	"net/http"
)

// GetCategories godoc
// @Summary Get all active categories
// @Description Retrieve all active support categories without authentication
// @Tags public
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /public/categories [get]
func GetCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Public categories endpoint - Coming soon"}`))
}