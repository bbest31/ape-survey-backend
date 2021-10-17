package main

import (
	"net/http"

	"github.com/apesurvey/ape-survey-backend/v2/routes"
	"github.com/apesurvey/ape-survey-backend/v2/server"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// routes
	router.HandleFunc("/users", routes.GetUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/users", routes.GetUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/users", routes.GetUserHandler).Methods(http.MethodDelete)
	router.HandleFunc("/users", routes.GetUserHandler).Methods(http.MethodPatch)

	router.HandleFunc("/reward-pools", routes.GetUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/reward-pools", routes.GetUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/reward-pools", routes.GetUserHandler).Methods(http.MethodDelete)
	router.HandleFunc("/reward-pools", routes.GetUserHandler).Methods(http.MethodPatch)

	router.HandleFunc("/rewards", routes.GetUserHandler).Methods(http.MethodGet)
	router.HandleFunc("/rewards", routes.GetUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/rewards", routes.GetUserHandler).Methods(http.MethodDelete)
	router.HandleFunc("/rewards", routes.GetUserHandler).Methods(http.MethodPatch)

	// add middleware
	router.Use(server.ValidateAccessToken)

}
