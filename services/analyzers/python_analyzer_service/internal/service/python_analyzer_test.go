package service

import (
	"testing"

	"python_analyzer_service/internal/models"
)

func TestPythonAnalyzerService_Analyze(t *testing.T) {
	service := NewPythonAnalyzerService()

	tests := []struct {
		name     string
		req      *models.AnalyzeRequest
		wantErr  bool
		checkLen bool
	}{
		{
			name: "empty files",
			req: &models.AnalyzeRequest{
				Files: []models.FileInput{},
			},
			wantErr:  false,
			checkLen: true,
		},
		{
			name: "non-python file",
			req: &models.AnalyzeRequest{
				Files: []models.FileInput{
					{Path: "test.txt", Content: "hello"},
				},
			},
			wantErr:  false,
			checkLen: true,
		},
		{
			name: "python file",
			req: &models.AnalyzeRequest{
				Files: []models.FileInput{
					{Path: "test.py", Content: "print('hello')"},
				},
			},
			wantErr:  false,
			checkLen: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Analyze(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp != nil {
				if tt.checkLen && len(resp.Files) != len(tt.req.Files) {
					t.Errorf("Analyze() returned %d files, want %d", len(resp.Files), len(tt.req.Files))
				}
			}
		})
	}
}

