package main

import (
	"log"
	"net/http"
	"os"

	"cpp_analyzer_service/internal/controller"
	"cpp_analyzer_service/internal/service"
)

func main() {
	analyzerService := service.NewCppAnalyzerService()
	analyzerController := controller.NewAnalyzerController(analyzerService)
	router := controller.NewRouter(analyzerController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	log.Printf("Starting cpp_analyzer_service on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

