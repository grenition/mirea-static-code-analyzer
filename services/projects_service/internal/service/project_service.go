package service

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"projects_service/internal/models"
	"projects_service/internal/repository"
)

type ProjectService struct {
	projectRepo    *repository.ProjectRepository
	fileRepo       *repository.FileRepository
	jwtSecret      string
	analyzerBaseURL string
}

func NewProjectService(projectRepo *repository.ProjectRepository, fileRepo *repository.FileRepository, jwtSecret, analyzerBaseURL string) *ProjectService {
	return &ProjectService{
		projectRepo:     projectRepo,
		fileRepo:        fileRepo,
		jwtSecret:       jwtSecret,
		analyzerBaseURL: analyzerBaseURL,
	}
}

func (s *ProjectService) ListProjects(token string) ([]*models.Project, error) {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	return s.projectRepo.FindByUserID(userID)
}

func (s *ProjectService) GetProject(token string, projectID int) (*models.ProjectWithFiles, error) {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, errors.New("forbidden")
	}

	files, err := s.fileRepo.FindByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	// Convert []*File to []File
	fileList := make([]models.File, len(files))
	for i, f := range files {
		fileList[i] = *f
	}

	return &models.ProjectWithFiles{
		Project: *project,
		Files:   fileList,
	}, nil
}

func (s *ProjectService) CreateProject(token string, req *models.CreateProjectRequest) (*models.Project, error) {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	if req.Name == "" {
		return nil, errors.New("project name is required")
	}

	project := &models.Project{
		Name:   req.Name,
		UserID: userID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(token string, projectID int, req *models.UpdateProjectRequest) (*models.Project, error) {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, errors.New("forbidden")
	}

	if req.Name == "" {
		return nil, errors.New("project name is required")
	}

	project.Name = req.Name
	if err := s.projectRepo.Update(project); err != nil {
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return project, nil
}

func (s *ProjectService) DeleteProject(token string, projectID int) error {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return errors.New("unauthorized")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	if project.UserID != userID {
		return errors.New("forbidden")
	}

	return s.projectRepo.Delete(projectID)
}

func (s *ProjectService) CreateFileFromZip(token string, projectID int, zipData []byte) error {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return errors.New("unauthorized")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return err
	}

	if project.UserID != userID {
		return errors.New("forbidden")
	}

	if len(zipData) > 25*1024*1024 {
		return errors.New("archive size exceeds 25 MB limit")
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to read zip: %w", err)
	}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			continue
		}

		// Normalize path
		path := strings.TrimPrefix(filepath.ToSlash(file.Name), "./")

		fileModel := &models.File{
			ProjectID: projectID,
			Path:      path,
			Content:   string(content),
		}

		s.fileRepo.Create(fileModel)
	}

	return nil
}

func (s *ProjectService) CreateFile(token string, projectID int, req *models.CreateFileRequest) (*models.File, error) {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	project, err := s.projectRepo.FindByID(projectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, errors.New("forbidden")
	}

	if req.Path == "" {
		return nil, errors.New("file path is required")
	}

	file := &models.File{
		ProjectID: projectID,
		Path:      req.Path,
		Content:   req.Content,
	}

	if err := s.fileRepo.Create(file); err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return file, nil
}

func (s *ProjectService) UpdateFile(token string, fileID int, req *models.UpdateFileRequest) (*models.File, error) {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.FindByID(file.ProjectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, errors.New("forbidden")
	}

	if req.Path == "" {
		return nil, errors.New("file path is required")
	}

	file.Path = req.Path
	file.Content = req.Content

	if err := s.fileRepo.Update(file); err != nil {
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	return file, nil
}

func (s *ProjectService) DeleteFile(token string, fileID int) error {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return errors.New("unauthorized")
	}

	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return err
	}

	project, err := s.projectRepo.FindByID(file.ProjectID)
	if err != nil {
		return err
	}

	if project.UserID != userID {
		return errors.New("forbidden")
	}

	return s.fileRepo.Delete(fileID)
}

func (s *ProjectService) AnalyzeFile(token string, fileID int, analyzerType string) (interface{}, error) {
	userID, _, err := validateToken(token, s.jwtSecret)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	file, err := s.fileRepo.FindByID(fileID)
	if err != nil {
		return nil, err
	}

	project, err := s.projectRepo.FindByID(file.ProjectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, errors.New("forbidden")
	}

	analyzerURL := fmt.Sprintf("%s/api/analyzer/%s", s.analyzerBaseURL, analyzerType)

	requestBody := map[string]interface{}{
		"files": []map[string]string{
			{
				"path":    file.Path,
				"content": file.Content,
			},
		},
	}

	jsonBody, _ := json.Marshal(requestBody)
	resp, err := http.Post(analyzerURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to call analyzer: %w", err)
	}
	defer resp.Body.Close()

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode analyzer response: %w", err)
	}

	return result, nil
}

