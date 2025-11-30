package main

import (
	"log"
	"net/http"
	"os"

	"projects_service/internal/config"
	"projects_service/internal/controller"
	"projects_service/internal/repository"
	"projects_service/internal/service"
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

	projectRepo := repository.NewProjectRepository(db)
	fileRepo := repository.NewFileRepository(db)
	projectService := service.NewProjectService(projectRepo, fileRepo, cfg.JWTSecret, cfg.AnalyzerBaseURL)
	projectController := controller.NewProjectController(projectService)

	router := controller.NewRouter(projectController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Starting projects_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

