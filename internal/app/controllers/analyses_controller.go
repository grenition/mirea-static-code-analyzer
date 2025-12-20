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
}

func NewAnalysesController(
	tmpl *template.Template,
	analysisRepo *repositories.AnalysisRepository,
	projectRepo *repositories.ProjectRepository,
	issueRepo *repositories.IssueRepository,
) *AnalysesController {
	return &AnalysesController{
		tmpl:         tmpl,
		analysisRepo: analysisRepo,
		projectRepo:  projectRepo,
		issueRepo:    issueRepo,
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

	var summary models.AnalysisSummary
	if analysis.SummaryJSON != nil {
		json.Unmarshal(analysis.SummaryJSON, &summary)
	}

	data := DetailsData{
		Analysis: analysis,
		Project:  project,
		Issues:   issues,
		Summary:  summary,
	}

	if err := executeTemplate(w, c.tmpl, "analysis_details", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

