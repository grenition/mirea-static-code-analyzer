package main

import (
	"log"
	"net/http"
	"os"

	"javascript_analyzer_service/internal/controller"
	"javascript_analyzer_service/internal/service"
)

func main() {
	analyzerService := service.NewJavaScriptAnalyzerService()
	analyzerController := controller.NewAnalyzerController(analyzerService)
	router := controller.NewRouter(analyzerController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Starting javascript_analyzer_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

