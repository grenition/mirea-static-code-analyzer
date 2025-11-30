package controller

import (
	"encoding/json"
	"net/http"

	"json_analyzer_service/internal/models"
	"json_analyzer_service/internal/service"
)

type AnalyzerController struct {
	service *service.JSONAnalyzerService
}

func NewAnalyzerController(service *service.JSONAnalyzerService) *AnalyzerController {
	return &AnalyzerController{service: service}
}

func (c *AnalyzerController) Analyze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AnalyzeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := c.service.Analyze(&req)
	if err != nil {
		respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, resp, http.StatusOK)
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

