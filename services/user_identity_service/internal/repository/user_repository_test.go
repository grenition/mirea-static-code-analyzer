package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"user_identity_service/internal/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://user:password@localhost:5432/test_user_identity_db?sslmode=disable"
	}

	db, err := NewPostgresDB(databaseURL)
	if err != nil {
		t.Skipf("Skipping test: failed to connect to test database: %v", err)
	}

	// Clean up
	db.Exec("DROP TABLE IF EXISTS users")
	RunMigrations(db)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user := &models.User{
		Username: "testuser",
		Password: "hashedpassword",
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if user.ID == 0 {
		t.Error("Create() should set user ID")
	}
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create a user first
	user := &models.User{
		Username: "finduser",
		Password: "hashedpassword",
	}
	repo.Create(user)

	// Find the user
	found, err := repo.FindByUsername("finduser")
	if err != nil {
		t.Fatalf("FindByUsername() error = %v", err)
	}

	if found.Username != "finduser" {
		t.Errorf("FindByUsername() username = %v, want finduser", found.Username)
	}

	// Test not found
	_, err = repo.FindByUsername("nonexistent")
	if err == nil {
		t.Error("FindByUsername() should return error for nonexistent user")
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create a user first
	user := &models.User{
		Username: "iduser",
		Password: "hashedpassword",
	}
	repo.Create(user)

	// Find the user
	found, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}

	if found.ID != user.ID {
		t.Errorf("FindByID() id = %v, want %v", found.ID, user.ID)
	}

	// Test not found
	_, err = repo.FindByID(99999)
	if err == nil {
		t.Error("FindByID() should return error for nonexistent user")
	}
}

