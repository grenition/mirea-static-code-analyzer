package service

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"javascript_analyzer_service/internal/models"
)

type JavaScriptAnalyzerService struct {
	eslintPath string
}

func NewJavaScriptAnalyzerService() *JavaScriptAnalyzerService {
	eslintPath := os.Getenv("ESLINT_PATH")
	if eslintPath == "" {
		eslintPath = "eslint"
	}

	return &JavaScriptAnalyzerService{
		eslintPath: eslintPath,
	}
}

func (s *JavaScriptAnalyzerService) Analyze(req *models.AnalyzeRequest) (*models.AnalyzeResponse, error) {
	var results []models.FileResult

	for _, file := range req.Files {
		if !strings.HasSuffix(strings.ToLower(file.Path), ".js") && !strings.HasSuffix(strings.ToLower(file.Path), ".jsx") {
			results = append(results, models.FileResult{
				Path:         file.Path,
				Comment:      "Not a JavaScript file",
				LineComments: []models.LineComment{},
			})
			continue
		}

		result := s.analyzeFile(file.Path, file.Content)
		results = append(results, result)
	}

	return &models.AnalyzeResponse{Files: results}, nil
}

func (s *JavaScriptAnalyzerService) analyzeFile(path, content string) models.FileResult {
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("analyze_%d_%s", time.Now().UnixNano(), filepath.Base(path)))
	defer os.Remove(tmpFile)

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		return models.FileResult{
			Path:         path,
			Comment:      fmt.Sprintf("Error: failed to create temp file: %v", err),
			LineComments: []models.LineComment{},
		}
	}

	cmd := exec.Command(s.eslintPath, "--format", "json", tmpFile)
	output, err := cmd.CombinedOutput()

	if err != nil && len(output) == 0 {
		return models.FileResult{
			Path:         path,
			Comment:      "OK",
			LineComments: []models.LineComment{},
		}
	}

	var eslintResults []map[string]interface{}
	if err := json.Unmarshal(output, &eslintResults); err != nil {
		return models.FileResult{
			Path:         path,
			Comment:      "OK",
			LineComments: []models.LineComment{},
		}
	}

	var lineComments []models.LineComment
	comment := "OK"

	for _, fileResult := range eslintResults {
		messages, _ := fileResult["messages"].([]interface{})
		for _, msg := range messages {
			msgMap, _ := msg.(map[string]interface{})
			line, _ := msgMap["line"].(float64)
			message, _ := msgMap["message"].(string)

			lineComments = append(lineComments, models.LineComment{
				Line:    int(line),
				Comment: message,
			})

			if comment == "OK" {
				comment = "Issues found"
			}
		}
	}

	return models.FileResult{
		Path:         path,
		Comment:      comment,
		LineComments: lineComments,
	}
}

