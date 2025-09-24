package handlers

import (
	"net/http"
)

// HealthCheck godoc
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"kemuko-support-service","platform":"kemuko"}`))
}