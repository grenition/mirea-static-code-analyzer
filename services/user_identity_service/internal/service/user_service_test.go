package service

import (
	"testing"

	"user_identity_service/internal/models"
)

func TestUserService_Register(t *testing.T) {
	// This is a unit test that would require a mock repository
	// For now, we'll test the validation logic
	service := NewUserService(nil, "test-secret")

	tests := []struct {
		name    string
		req     *models.RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty username",
			req:     &models.RegisterRequest{Username: "", Password: "password123"},
			wantErr: true,
			errMsg:  "username and password are required",
		},
		{
			name:    "empty password",
			req:     &models.RegisterRequest{Username: "user", Password: ""},
			wantErr: true,
			errMsg:  "username and password are required",
		},
		{
			name:    "short username",
			req:     &models.RegisterRequest{Username: "ab", Password: "password123"},
			wantErr: true,
			errMsg:  "username must be at least 3 characters",
		},
		{
			name:    "short password",
			req:     &models.RegisterRequest{Username: "user123", Password: "pass"},
			wantErr: true,
			errMsg:  "password must be at least 6 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Register() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	service := NewUserService(nil, "test-secret")

	tests := []struct {
		name    string
		req     *models.LoginRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "empty username",
			req:     &models.LoginRequest{Username: "", Password: "password123"},
			wantErr: true,
			errMsg:  "username and password are required",
		},
		{
			name:    "empty password",
			req:     &models.LoginRequest{Username: "user", Password: ""},
			wantErr: true,
			errMsg:  "username and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMsg {
				t.Errorf("Login() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUserService_ValidateToken(t *testing.T) {
	service := NewUserService(nil, "test-secret")

	// Generate a valid token
	userID := 1
	username := "testuser"
	token, err := service.generateToken(userID, username)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Test valid token
	validatedID, validatedUsername, err := service.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() error = %v", err)
	}
	if validatedID != userID {
		t.Errorf("ValidateToken() userID = %v, want %v", validatedID, userID)
	}
	if validatedUsername != username {
		t.Errorf("ValidateToken() username = %v, want %v", validatedUsername, username)
	}

	// Test invalid token
	_, _, err = service.ValidateToken("invalid-token")
	if err == nil {
		t.Error("ValidateToken() should return error for invalid token")
	}
}

