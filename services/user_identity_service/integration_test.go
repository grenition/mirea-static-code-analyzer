package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"user_identity_service/internal/config"
	"user_identity_service/internal/controller"
	"user_identity_service/internal/models"
	"user_identity_service/internal/repository"
	"user_identity_service/internal/service"
)

func setupIntegrationTest(t *testing.T) (*httptest.Server, func()) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/test_user_identity_db?sslmode=disable"
	}

	db, err := repository.NewPostgresDB(databaseURL)
	if err != nil {
		t.Skipf("Skipping integration test: failed to connect to test database: %v", err)
	}

	// Clean up
	db.Exec("DROP TABLE IF EXISTS users")
	repository.RunMigrations(db)

	cfg := config.Load()
	cfg.DatabaseURL = databaseURL

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	userController := controller.NewUserController(userService)
	router := controller.NewRouter(userController)

	server := httptest.NewServer(router)

	cleanup := func() {
		server.Close()
		db.Exec("DROP TABLE IF EXISTS users")
		db.Close()
	}

	return server, cleanup
}

func TestIntegration_Register(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	reqBody := models.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(reqBody)
	resp, err := http.Post(server.URL+"/api/users/register", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Register() status = %v, want %v", resp.StatusCode, http.StatusCreated)
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if authResp.Token == "" {
		t.Error("Register() should return a token")
	}
}

func TestIntegration_Login(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// First register
	registerReq := models.RegisterRequest{
		Username: "loginuser",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(registerReq)
	http.Post(server.URL+"/api/users/register", "application/json", bytes.NewBuffer(jsonBody))

	// Then login
	loginReq := models.LoginRequest{
		Username: "loginuser",
		Password: "password123",
	}
	jsonBody, _ = json.Marshal(loginReq)
	resp, err := http.Post(server.URL+"/api/users/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Login() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}

	var authResp models.AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if authResp.Token == "" {
		t.Error("Login() should return a token")
	}
}

func TestIntegration_LoginInvalidCredentials(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	loginReq := models.LoginRequest{
		Username: "nonexistent",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(loginReq)
	resp, err := http.Post(server.URL+"/api/users/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Login() status = %v, want %v", resp.StatusCode, http.StatusUnauthorized)
	}
}

