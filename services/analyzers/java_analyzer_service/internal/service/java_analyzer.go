package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"java_analyzer_service/internal/models"
)

type JavaAnalyzerService struct {
	checkstylePath string
}

func NewJavaAnalyzerService() *JavaAnalyzerService {
	checkstylePath := os.Getenv("CHECKSTYLE_PATH")
	if checkstylePath == "" {
		checkstylePath = "checkstyle"
	}

	return &JavaAnalyzerService{
		checkstylePath: checkstylePath,
	}
}

func (s *JavaAnalyzerService) Analyze(req *models.AnalyzeRequest) (*models.AnalyzeResponse, error) {
	var results []models.FileResult

	for _, file := range req.Files {
		if !strings.HasSuffix(strings.ToLower(file.Path), ".java") {
			results = append(results, models.FileResult{
				Path:         file.Path,
				Comment:      "Not a Java file",
				LineComments: []models.LineComment{},
			})
			continue
		}

		result := s.analyzeFile(file.Path, file.Content)
		results = append(results, result)
	}

	return &models.AnalyzeResponse{Files: results}, nil
}

func (s *JavaAnalyzerService) analyzeFile(path, content string) models.FileResult {
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

	cmd := exec.Command(s.checkstylePath, "-f", "plain", tmpFile)
	output, err := cmd.CombinedOutput()

	if err != nil && len(output) == 0 {
		return models.FileResult{
			Path:         path,
			Comment:      "OK",
			LineComments: []models.LineComment{},
		}
	}

	return s.parseOutput(path, string(output))
}

func (s *JavaAnalyzerService) parseOutput(path, output string) models.FileResult {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var lineComments []models.LineComment
	comment := "OK"

	for _, line := range lines {
		if line == "" {
			continue
		}

		// Parse checkstyle output format
		parts := strings.Split(line, ":")
		if len(parts) >= 3 {
			var lineNum int
			fmt.Sscanf(parts[1], "%d", &lineNum)
			msg := strings.Join(parts[2:], ":")

			lineComments = append(lineComments, models.LineComment{
				Line:    lineNum,
				Comment: strings.TrimSpace(msg),
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

