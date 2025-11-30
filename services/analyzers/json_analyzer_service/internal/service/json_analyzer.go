package service

import (
	"encoding/json"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"json_analyzer_service/internal/models"
)

type JSONAnalyzerService struct{}

func NewJSONAnalyzerService() *JSONAnalyzerService {
	return &JSONAnalyzerService{}
}

func (s *JSONAnalyzerService) Analyze(req *models.AnalyzeRequest) (*models.AnalyzeResponse, error) {
	var results []models.FileResult

	for _, file := range req.Files {
		if !strings.HasSuffix(strings.ToLower(file.Path), ".json") {
			results = append(results, models.FileResult{
				Path:         file.Path,
				Comment:      "Not a JSON file",
				LineComments: []models.LineComment{},
			})
			continue
		}

		result := s.analyzeFile(file.Path, file.Content)
		results = append(results, result)
	}

	return &models.AnalyzeResponse{Files: results}, nil
}

func (s *JSONAnalyzerService) analyzeFile(path, content string) models.FileResult {
	// First, validate JSON syntax
	var jsonData interface{}
	if err := json.Unmarshal([]byte(content), &jsonData); err != nil {
		line := s.getErrorLine(content, err.Error())
		return models.FileResult{
			Path:    path,
			Comment: "Invalid JSON syntax",
			LineComments: []models.LineComment{
				{
					Line:    line,
					Comment: err.Error(),
				},
			},
		}
	}

	// If schema validation is needed, it can be added here
	// For now, we just validate syntax
	schemaLoader := gojsonschema.NewStringLoader(`{}`)
	documentLoader := gojsonschema.NewStringLoader(content)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return models.FileResult{
			Path:         path,
			Comment:      "OK",
			LineComments: []models.LineComment{},
		}
	}

	if !result.Valid() {
		var lineComments []models.LineComment
		for _, desc := range result.Errors() {
			lineComments = append(lineComments, models.LineComment{
				Line:    1, // gojsonschema doesn't provide line numbers directly
				Comment: desc.String(),
			})
		}
		return models.FileResult{
			Path:         path,
			Comment:      "JSON validation issues found",
			LineComments: lineComments,
		}
	}

	return models.FileResult{
		Path:         path,
		Comment:      "OK",
		LineComments: []models.LineComment{},
	}
}

func (s *JSONAnalyzerService) getErrorLine(content, errorMsg string) int {
	// Simple heuristic to find line number from error
	// JSON unmarshal errors typically include position info
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		return len(lines)
	}
	return 1
}

