package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"cpp_analyzer_service/internal/models"
)

type CppAnalyzerService struct {
	cppcheckPath string
}

func NewCppAnalyzerService() *CppAnalyzerService {
	cppcheckPath := os.Getenv("CPPCHECK_PATH")
	if cppcheckPath == "" {
		cppcheckPath = "cppcheck"
	}

	return &CppAnalyzerService{
		cppcheckPath: cppcheckPath,
	}
}

func (s *CppAnalyzerService) Analyze(req *models.AnalyzeRequest) (*models.AnalyzeResponse, error) {
	var results []models.FileResult

	for _, file := range req.Files {
		ext := strings.ToLower(filepath.Ext(file.Path))
		if ext != ".cpp" && ext != ".c" && ext != ".cc" && ext != ".cxx" && ext != ".h" && ext != ".hpp" {
			results = append(results, models.FileResult{
				Path:         file.Path,
				Comment:      "Not a C/C++ file",
				LineComments: []models.LineComment{},
			})
			continue
		}

		result := s.analyzeFile(file.Path, file.Content)
		results = append(results, result)
	}

	return &models.AnalyzeResponse{Files: results}, nil
}

func (s *CppAnalyzerService) analyzeFile(path, content string) models.FileResult {
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

	cmd := exec.Command(s.cppcheckPath, "--enable=all", "--xml", tmpFile)
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

func (s *CppAnalyzerService) parseOutput(path, output string) models.FileResult {
	// Simple XML parsing for cppcheck output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	var lineComments []models.LineComment
	comment := "OK"

	for _, line := range lines {
		if strings.Contains(line, "<error") {
			// Extract line number and message from XML
			var lineNum int
			fmt.Sscanf(line, "%*[^l]line=\"%d\"", &lineNum)
			if lineNum > 0 {
				msg := extractMessage(line)
				lineComments = append(lineComments, models.LineComment{
					Line:    lineNum,
					Comment: msg,
				})
				if comment == "OK" {
					comment = "Issues found"
				}
			}
		}
	}

	return models.FileResult{
		Path:         path,
		Comment:      comment,
		LineComments: lineComments,
	}
}

func extractMessage(line string) string {
	start := strings.Index(line, "msg=\"")
	if start == -1 {
		return "Issue found"
	}
	start += 5
	end := strings.Index(line[start:], "\"")
	if end == -1 {
		return line[start:]
	}
	return line[start : start+end]
}

