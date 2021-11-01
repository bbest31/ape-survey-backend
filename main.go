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
	router.HandleFunc("/users", routes.DeleteUserHandler).Methods(http.MethodDelete)
	router.HandleFunc("/users", routes.PatchUserHandler).Methods(http.MethodPatch)

	router.HandleFunc("/reward-pools", routes.GetRewardPoolHandler).Methods(http.MethodGet)
	router.HandleFunc("/reward-pools", routes.PostRewardPoolHandler).Methods(http.MethodPost)
	router.HandleFunc("/reward-pools", routes.DeleteRewardPoolHandler).Methods(http.MethodDelete)
	router.HandleFunc("/reward-pools", routes.PatchRewardPoolHandler).Methods(http.MethodPatch)

	router.HandleFunc("/rewards", routes.GetRewardHandler).Methods(http.MethodGet)
	router.HandleFunc("/rewards", routes.PostRewardHandler).Methods(http.MethodPost)
	router.HandleFunc("/rewards", routes.DeleteRewardHandler).Methods(http.MethodDelete)

	// TODO look up survey monkey survey response webhook details.
	router.HandleFunc("/survey-response", routes.SurveyResponseWebhook).Methods(http.MethodPost)
	router.HandleFunc("/surveys/{id}", routes.GetUserSurveys).Methods(http.MethodGet)
	router.HandleFunc("/connect-surveymonkey", routes.ConnectSurveyMonkey).Methods(http.MethodPost)

	// add middleware
	router.Use(server.ValidateAccessToken)

}
