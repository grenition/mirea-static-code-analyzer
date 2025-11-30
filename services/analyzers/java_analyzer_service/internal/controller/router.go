package controller

import (
	"github.com/gorilla/mux"
)

func NewRouter(analyzerController *AnalyzerController) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/analyzer/java", analyzerController.Analyze).Methods("POST")
	return router
}

