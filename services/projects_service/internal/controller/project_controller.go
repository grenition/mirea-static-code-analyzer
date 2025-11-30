package controller

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"projects_service/internal/models"
	"projects_service/internal/service"
)

type ProjectController struct {
	service *service.ProjectService
}

func NewProjectController(service *service.ProjectService) *ProjectController {
	return &ProjectController{service: service}
}

func (c *ProjectController) getToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

func (c *ProjectController) ListProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := c.getToken(r)
	projects, err := c.service.ListProjects(token)
	if err != nil {
		respondError(w, err.Error(), http.StatusUnauthorized)
		return
	}

	respondJSON(w, projects, http.StatusOK)
}

func (c *ProjectController) GetProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	project, err := c.service.GetProject(token, projectID)
	if err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	respondJSON(w, project, http.StatusOK)
}

func (c *ProjectController) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	project, err := c.service.CreateProject(token, &req)
	if err != nil {
		if err.Error() == "unauthorized" {
			respondError(w, err.Error(), http.StatusUnauthorized)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	respondJSON(w, project, http.StatusCreated)
}

func (c *ProjectController) UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	project, err := c.service.UpdateProject(token, projectID, &req)
	if err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	respondJSON(w, project, http.StatusOK)
}

func (c *ProjectController) DeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	if err := c.service.DeleteProject(token, projectID); err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ProjectController) CreateProjectFromZip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	zipData, err := io.ReadAll(r.Body)
	if err != nil {
		respondError(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	if err := c.service.CreateFileFromZip(token, projectID, zipData); err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (c *ProjectController) ListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	project, err := c.service.GetProject(token, projectID)
	if err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	respondJSON(w, project.Files, http.StatusOK)
}

func (c *ProjectController) CreateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondError(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req models.CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	file, err := c.service.CreateFile(token, projectID, &req)
	if err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	respondJSON(w, file, http.StatusCreated)
}

func (c *ProjectController) UpdateFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["fileId"])
	if err != nil {
		respondError(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	file, err := c.service.UpdateFile(token, fileID, &req)
	if err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	respondJSON(w, file, http.StatusOK)
}

func (c *ProjectController) DeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["fileId"])
	if err != nil {
		respondError(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	if err := c.service.DeleteFile(token, fileID); err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusNotFound)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *ProjectController) AnalyzeFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	fileID, err := strconv.Atoi(vars["fileId"])
	if err != nil {
		respondError(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	analyzerType := r.URL.Query().Get("analyzer")
	if analyzerType == "" {
		respondError(w, "analyzer parameter is required", http.StatusBadRequest)
		return
	}

	token := c.getToken(r)
	result, err := c.service.AnalyzeFile(token, fileID, analyzerType)
	if err != nil {
		if err.Error() == "unauthorized" || err.Error() == "forbidden" {
			respondError(w, err.Error(), http.StatusForbidden)
		} else {
			respondError(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	respondJSON(w, result, http.StatusOK)
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

