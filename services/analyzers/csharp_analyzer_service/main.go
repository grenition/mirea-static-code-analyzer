package main

import (
	"log"
	"net/http"
	"os"

	"csharp_analyzer_service/internal/controller"
	"csharp_analyzer_service/internal/service"
)

func main() {
	analyzerService := service.NewCsharpAnalyzerService()
	analyzerController := controller.NewAnalyzerController(analyzerService)
	router := controller.NewRouter(analyzerController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}

	log.Printf("Starting csharp_analyzer_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

