package service

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"python_analyzer_service/internal/models"
)

type PythonAnalyzerService struct {
	flake8Path string
}

func NewPythonAnalyzerService() *PythonAnalyzerService {
	flake8Path := os.Getenv("FLAKE8_PATH")
	if flake8Path == "" {
		flake8Path = "flake8"
	}

	return &PythonAnalyzerService{
		flake8Path: flake8Path,
	}
}

func (s *PythonAnalyzerService) Analyze(req *models.AnalyzeRequest) (*models.AnalyzeResponse, error) {
	var results []models.FileResult

	for _, file := range req.Files {
		if !strings.HasSuffix(strings.ToLower(file.Path), ".py") {
			results = append(results, models.FileResult{
				Path:         file.Path,
				Comment:      "Not a Python file",
				LineComments: []models.LineComment{},
			})
			continue
		}

		result := s.analyzeFile(file.Path, file.Content)
		results = append(results, result)
	}

	return &models.AnalyzeResponse{Files: results}, nil
}

func (s *PythonAnalyzerService) analyzeFile(path, content string) models.FileResult {
	// Create temporary file
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

	// Run flake8
	cmd := exec.Command(s.flake8Path, "--format=json", tmpFile)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// flake8 returns non-zero exit code if issues found, so we need to parse output
		if len(output) == 0 {
			return models.FileResult{
				Path:         path,
				Comment:      "OK",
				LineComments: []models.LineComment{},
			}
		}
	}

	// Parse flake8 JSON output
	var flake8Results []map[string]interface{}
	if err := json.Unmarshal(output, &flake8Results); err != nil {
		// If JSON parsing fails, try to parse text output
		return s.parseTextOutput(path, string(output))
	}

	var lineComments []models.LineComment
	comment := "OK"

	for _, result := range flake8Results {
		lineNum, _ := result["line_number"].(float64)
		code, _ := result["code"].(string)
		text, _ := result["text"].(string)

		lineComments = append(lineComments, models.LineComment{
			Line:    int(lineNum),
			Comment: fmt.Sprintf("%s: %s", code, text),
		})

		if comment == "OK" {
			comment = "Issues found"
		}
	}

	return models.FileResult{
		Path:         path,
		Comment:      comment,
		LineComments: lineComments,
	}
}

func (s *PythonAnalyzerService) parseTextOutput(path, output string) models.FileResult {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var lineComments []models.LineComment
	comment := "OK"

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse flake8 text format: path:line:col: code message
		parts := strings.SplitN(line, ":", 4)
		if len(parts) >= 3 {
			var lineNum int
			fmt.Sscanf(parts[1], "%d", &lineNum)
			msg := ""
			if len(parts) >= 4 {
				msg = strings.TrimSpace(parts[3])
			}

			lineComments = append(lineComments, models.LineComment{
				Line:    lineNum,
				Comment: msg,
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

