package main

import (
	"log"
	"net/http"
	"os"

	"java_analyzer_service/internal/controller"
	"java_analyzer_service/internal/service"
)

func main() {
	analyzerService := service.NewJavaAnalyzerService()
	analyzerController := controller.NewAnalyzerController(analyzerService)
	router := controller.NewRouter(analyzerController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	log.Printf("Starting java_analyzer_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

