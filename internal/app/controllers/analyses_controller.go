package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"webapp/internal/app/models"
	"webapp/internal/app/repositories"
)

type AnalysesController struct {
	tmpl         *template.Template
	analysisRepo *repositories.AnalysisRepository
	projectRepo  *repositories.ProjectRepository
	issueRepo    *repositories.IssueRepository
	fileRepo     *repositories.FileRepository
}

func NewAnalysesController(
	tmpl *template.Template,
	analysisRepo *repositories.AnalysisRepository,
	projectRepo *repositories.ProjectRepository,
	issueRepo *repositories.IssueRepository,
	fileRepo *repositories.FileRepository,
) *AnalysesController {
	return &AnalysesController{
		tmpl:         tmpl,
		analysisRepo: analysisRepo,
		projectRepo:  projectRepo,
		issueRepo:    issueRepo,
		fileRepo:     fileRepo,
	}
}

type ListData struct {
	Analyses []AnalysisWithProject
}

type AnalysisWithProject struct {
	Analysis *models.AnalysisRun
	Project  *models.Project
}

type DetailsData struct {
	Analysis *models.AnalysisRun
	Project  *models.Project
	Issues   []*models.Issue
	Summary  models.AnalysisSummary
	Files    []*models.AnalysisFile
}

func (c *AnalysesController) List(w http.ResponseWriter, r *http.Request) {
	analyses, err := c.analysisRepo.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var data ListData
	for _, analysis := range analyses {
		project, err := c.projectRepo.GetByID(analysis.ProjectID)
		if err != nil {
			continue
		}
	data.Analyses = append(data.Analyses, AnalysisWithProject{
		Analysis: analysis,
		Project:  project,
	})
}

	if err := executeTemplate(w, c.tmpl, "analyses", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AnalysesController) Details(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid analysis ID", http.StatusBadRequest)
		return
	}

	analysis, err := c.analysisRepo.GetByID(id)
	if err != nil {
		http.Error(w, "Analysis not found", http.StatusNotFound)
		return
	}

	project, err := c.projectRepo.GetByID(analysis.ProjectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	issues, err := c.issueRepo.GetByAnalysisID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	files, err := c.fileRepo.GetByAnalysisID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var summary models.AnalysisSummary
	if analysis.SummaryJSON != nil {
		json.Unmarshal(analysis.SummaryJSON, &summary)
	}

	data := DetailsData{
		Analysis: analysis,
		Project:  project,
		Issues:   issues,
		Summary:  summary,
		Files:    files,
	}

	if err := executeTemplate(w, c.tmpl, "analysis_details", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *AnalysesController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid analysis ID", http.StatusBadRequest)
		return
	}

	// Delete analysis (cascade will delete issues and files)
	if err := c.analysisRepo.Delete(id); err != nil {
		http.Error(w, "Failed to delete analysis: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/analyses", http.StatusSeeOther)
}

func (c *AnalysesController) GetFileContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID, err := uuid.Parse(vars["fileId"])
	if err != nil {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, err := c.fileRepo.GetByID(fileID)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(file.Content)
}

func (c *AnalysesController) GetFileByPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	analysisID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "Invalid analysis ID", http.StatusBadRequest)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	file, err := c.fileRepo.GetByAnalysisIDAndPath(analysisID, filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(file.Content)
}

