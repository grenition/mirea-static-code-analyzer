package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"webapp/internal/app/models"
	"webapp/internal/app/repositories"
	"webapp/internal/app/services"
)

type UploadController struct {
	tmpl            *template.Template
	analyzerService *services.AnalyzerService
	storageService  *services.StorageService
	projectRepo     *repositories.ProjectRepository
	analysisRepo    *repositories.AnalysisRepository
	issueRepo       *repositories.IssueRepository
}

func NewUploadController(
	tmpl *template.Template,
	analyzerService *services.AnalyzerService,
	storageService *services.StorageService,
	projectRepo *repositories.ProjectRepository,
	analysisRepo *repositories.AnalysisRepository,
	issueRepo *repositories.IssueRepository,
) *UploadController {
	return &UploadController{
		tmpl:            tmpl,
		analyzerService: analyzerService,
		storageService:  storageService,
		projectRepo:     projectRepo,
		analysisRepo:    analysisRepo,
		issueRepo:       issueRepo,
	}
}

func (c *UploadController) ShowUpload(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, c.tmpl, "upload", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *UploadController) ShowSnippetAnalyzer(w http.ResponseWriter, r *http.Request) {
	if err := executeTemplate(w, c.tmpl, "snippet_analyzer", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *UploadController) HandleZipUpload(w http.ResponseWriter, r *http.Request) {
	const maxSize = 20 * 1024 * 1024 // 20 MB
	r.ParseMultipartForm(maxSize)

	projectName := r.FormValue("project_name")
	if projectName == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("zip_file")
	if err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > maxSize {
		http.Error(w, "File too large (max 20MB)", http.StatusBadRequest)
		return
	}

	// Create project
	project, err := c.projectRepo.Create(projectName)
	if err != nil {
		http.Error(w, "Failed to create project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save uploaded file temporarily
	tmpFile, err := os.CreateTemp("", "upload-*.zip")
	if err != nil {
		http.Error(w, "Failed to create temp file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create analysis run
	inputMeta, _ := json.Marshal(map[string]interface{}{
		"file_name": header.Filename,
		"size_bytes": header.Size,
	})

	analysis, err := c.analysisRepo.Create(project.ID, "zip", inputMeta)
	if err != nil {
		http.Error(w, "Failed to create analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update status to Running
	c.analysisRepo.UpdateStatus(analysis.ID, "Running")

	// Extract and analyze
	zipPath, err := c.storageService.SaveZipFile(tmpFile.Name(), analysis.ID.String())
	if err != nil {
		c.analysisRepo.UpdateStatus(analysis.ID, "Failed")
		http.Error(w, "Failed to save zip: "+err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := c.storageService.ExtractZip(zipPath, "")
	if err != nil {
		c.analysisRepo.UpdateStatus(analysis.ID, "Failed")
		http.Error(w, "Failed to extract zip: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(files) == 0 {
		c.analysisRepo.UpdateStatus(analysis.ID, "Failed")
		http.Error(w, "No valid files found in zip", http.StatusBadRequest)
		return
	}

	// Run analysis
	issues, summary := c.analyzerService.Analyze(files, analysis.ID.String())

	// Save summary
	summaryJSON, _ := json.Marshal(summary)
	c.analysisRepo.UpdateSummary(analysis.ID, summaryJSON)

	// Save issues
	issueModels := c.analyzerService.ConvertToModels(issues, analysis.ID)
	if err := c.issueRepo.CreateBatch(issueModels); err != nil {
		c.analysisRepo.UpdateStatus(analysis.ID, "Failed")
		http.Error(w, "Failed to save issues: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Mark as done
	c.analysisRepo.UpdateStatus(analysis.ID, "Done")

	http.Redirect(w, r, fmt.Sprintf("/analyses/%s", analysis.ID), http.StatusSeeOther)
}

func (c *UploadController) HandleSnippetUpload(w http.ResponseWriter, r *http.Request) {
	const maxSize = 50000 // 50,000 characters

	projectName := r.FormValue("project_name")
	if projectName == "" {
		http.Error(w, "Project name is required", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	if len(code) > maxSize {
		http.Error(w, "Code too long (max 50,000 characters)", http.StatusBadRequest)
		return
	}

	ext := r.FormValue("extension")
	if ext == "" {
		ext = ".txt"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// Create project
	project, err := c.projectRepo.Create(projectName)
	if err != nil {
		http.Error(w, "Failed to create project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create analysis run
	inputMeta, _ := json.Marshal(map[string]interface{}{
		"extension": ext,
		"size_bytes": len(code),
	})

	analysis, err := c.analysisRepo.Create(project.ID, "snippet", inputMeta)
	if err != nil {
		http.Error(w, "Failed to create analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update status to Running
	c.analysisRepo.UpdateStatus(analysis.ID, "Running")

	// Create single file input
	filePath := "snippet" + ext
	files := []services.FileInput{
		{
			Path:    filePath,
			Content: []byte(code),
		},
	}

	// Run analysis
	issues, summary := c.analyzerService.Analyze(files, analysis.ID.String())

	// Save summary
	summaryJSON, _ := json.Marshal(summary)
	c.analysisRepo.UpdateSummary(analysis.ID, summaryJSON)

	// Save issues
	issueModels := c.analyzerService.ConvertToModels(issues, analysis.ID)
	if err := c.issueRepo.CreateBatch(issueModels); err != nil {
		c.analysisRepo.UpdateStatus(analysis.ID, "Failed")
		http.Error(w, "Failed to save issues: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Mark as done
	c.analysisRepo.UpdateStatus(analysis.ID, "Done")

	http.Redirect(w, r, fmt.Sprintf("/analyses/%s", analysis.ID), http.StatusSeeOther)
}

// HandleSnippetAnalyzeAPI handles API requests for live code analysis (returns JSON)
func (c *UploadController) HandleSnippetAnalyzeAPI(w http.ResponseWriter, r *http.Request) {
	const maxSize = 50000 // 50,000 characters

	w.Header().Set("Content-Type", "application/json")

	// Use default project name for live analyzer
	projectName := "Live Analysis"

	code := r.FormValue("code")
	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Code is required"})
		return
	}

	if len(code) > maxSize {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("Code too long (max %d characters)", maxSize)})
		return
	}

	ext := r.FormValue("extension")
	if ext == "" {
		ext = ".txt"
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// Create project (optional for API - we can skip saving if needed, but let's save it)
	project, err := c.projectRepo.Create(projectName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create project: " + err.Error()})
		return
	}

	// Create analysis run
	inputMeta, _ := json.Marshal(map[string]interface{}{
		"extension": ext,
		"size_bytes": len(code),
	})

	analysis, err := c.analysisRepo.Create(project.ID, "snippet", inputMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create analysis: " + err.Error()})
		return
	}

	// Update status to Running
	c.analysisRepo.UpdateStatus(analysis.ID, "Running")

	// Create single file input
	filePath := "snippet" + ext
	files := []services.FileInput{
		{
			Path:    filePath,
			Content: []byte(code),
		},
	}

	// Run analysis
	issues, summary := c.analyzerService.Analyze(files, analysis.ID.String())

	// Save summary
	summaryJSON, _ := json.Marshal(summary)
	c.analysisRepo.UpdateSummary(analysis.ID, summaryJSON)

	// Save issues
	issueModels := c.analyzerService.ConvertToModels(issues, analysis.ID)
	if err := c.issueRepo.CreateBatch(issueModels); err != nil {
		c.analysisRepo.UpdateStatus(analysis.ID, "Failed")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to save issues: " + err.Error()})
		return
	}

	// Mark as done
	c.analysisRepo.UpdateStatus(analysis.ID, "Done")

	// Return JSON response
	response := map[string]interface{}{
		"analysis_id": analysis.ID.String(),
		"summary":     summary,
		"issues":      convertIssuesToJSON(issueModels),
	}

	json.NewEncoder(w).Encode(response)
}

func convertIssuesToJSON(issues []*models.Issue) []map[string]interface{} {
	result := make([]map[string]interface{}, len(issues))
	for i, issue := range issues {
		result[i] = map[string]interface{}{
			"severity": issue.Severity,
			"rule_code": issue.RuleCode,
			"message":   issue.Message,
			"file_path": issue.FilePath,
			"line":      issue.Line,
		}
	}
	return result
}

