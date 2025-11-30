package main

import (
	"log"
	"net/http"
	"os"

	"python_analyzer_service/internal/controller"
	"python_analyzer_service/internal/service"
)

func main() {
	analyzerService := service.NewPythonAnalyzerService()
	analyzerController := controller.NewAnalyzerController(analyzerService)
	router := controller.NewRouter(analyzerController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting python_analyzer_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

