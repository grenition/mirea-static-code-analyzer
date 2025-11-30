package main

import (
	"log"
	"net/http"
	"os"

	"json_analyzer_service/internal/controller"
	"json_analyzer_service/internal/service"
)

func main() {
	analyzerService := service.NewJSONAnalyzerService()
	analyzerController := controller.NewAnalyzerController(analyzerService)
	router := controller.NewRouter(analyzerController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	log.Printf("Starting json_analyzer_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

