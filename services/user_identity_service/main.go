package main

import (
	"log"
	"net/http"
	"os"

	"user_identity_service/internal/config"
	"user_identity_service/internal/controller"
	"user_identity_service/internal/repository"
	"user_identity_service/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := repository.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	userController := controller.NewUserController(userService)

	router := controller.NewRouter(userController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting user_identity_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

