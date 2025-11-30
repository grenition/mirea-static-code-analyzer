package controller

import (
	"github.com/gorilla/mux"
)

func NewRouter(userController *UserController) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/users/register", userController.Register).Methods("POST")
	router.HandleFunc("/api/users/login", userController.Login).Methods("POST")

	return router
}

