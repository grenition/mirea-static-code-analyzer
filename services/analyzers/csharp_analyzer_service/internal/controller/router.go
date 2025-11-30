package controller

import (
	"github.com/gorilla/mux"
)

func NewRouter(analyzerController *AnalyzerController) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/analyzer/csharp", analyzerController.Analyze).Methods("POST")
	return router
}

