package service

import (
	"testing"

	"projects_service/internal/models"
)

func TestProjectService_ValidateToken(t *testing.T) {
	service := NewProjectService(nil, nil, "test-secret", "http://localhost")

	// Test with invalid token
	_, _, err := validateToken("invalid-token", "test-secret")
	if err == nil {
		t.Error("validateToken() should return error for invalid token")
	}
}

