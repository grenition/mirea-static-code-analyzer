package service

import (
	"testing"

	"json_analyzer_service/internal/models"
)

func TestJSONAnalyzerService_Analyze(t *testing.T) {
	service := NewJSONAnalyzerService()

	tests := []struct {
		name    string
		req     *models.AnalyzeRequest
		wantErr bool
	}{
		{
			name: "valid JSON",
			req: &models.AnalyzeRequest{
				Files: []models.FileInput{
					{Path: "test.json", Content: `{"key": "value"}`},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid JSON",
			req: &models.AnalyzeRequest{
				Files: []models.FileInput{
					{Path: "test.json", Content: `{"key": "value"`},
				},
			},
			wantErr: false,
		},
		{
			name: "non-JSON file",
			req: &models.AnalyzeRequest{
				Files: []models.FileInput{
					{Path: "test.txt", Content: "hello"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Analyze(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Analyze() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if resp == nil {
				t.Error("Analyze() returned nil response")
			}
		})
	}
}

