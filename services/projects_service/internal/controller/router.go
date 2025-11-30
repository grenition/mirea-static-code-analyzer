package controller

import (
	"github.com/gorilla/mux"
)

func NewRouter(projectController *ProjectController) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/projects", projectController.ListProjects).Methods("GET")
	router.HandleFunc("/api/projects", projectController.CreateProject).Methods("POST")
	router.HandleFunc("/api/projects/{id}", projectController.GetProject).Methods("GET")
	router.HandleFunc("/api/projects/{id}", projectController.UpdateProject).Methods("PUT")
	router.HandleFunc("/api/projects/{id}", projectController.DeleteProject).Methods("DELETE")
	router.HandleFunc("/api/projects/{id}/upload", projectController.CreateProjectFromZip).Methods("POST")
	router.HandleFunc("/api/projects/{id}/files", projectController.ListFiles).Methods("GET")
	router.HandleFunc("/api/projects/{id}/files", projectController.CreateFile).Methods("POST")
	router.HandleFunc("/api/projects/{id}/files/{fileId}", projectController.UpdateFile).Methods("PUT")
	router.HandleFunc("/api/projects/{id}/files/{fileId}", projectController.DeleteFile).Methods("DELETE")
	router.HandleFunc("/api/projects/{id}/files/{fileId}/analyze", projectController.AnalyzeFile).Methods("POST")

	return router
}

