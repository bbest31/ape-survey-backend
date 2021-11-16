package main

import (
	"log"
	"net/http"

	"github.com/apesurvey/ape-survey-backend/v2/routes"
	"github.com/apesurvey/ape-survey-backend/v2/server"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// routes
	router.HandleFunc("/user/{id}", routes.DeleteUserHandler).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/user/{id}", routes.PatchUserHandler).Methods(http.MethodPatch)

	router.HandleFunc("/reward-pools", routes.GetRewardPoolHandler).Methods(http.MethodGet)
	router.HandleFunc("/reward-pools", routes.PostRewardPoolHandler).Methods(http.MethodPost)
	router.HandleFunc("/reward-pools", routes.DeleteRewardPoolHandler).Methods(http.MethodDelete)
	router.HandleFunc("/reward-pools", routes.PatchRewardPoolHandler).Methods(http.MethodPatch)

	router.HandleFunc("/rewards", routes.GetRewardHandler).Methods(http.MethodGet)
	router.HandleFunc("/rewards", routes.PostRewardHandler).Methods(http.MethodPost)
	router.HandleFunc("/rewards", routes.DeleteRewardHandler).Methods(http.MethodDelete)

	// TODO look up survey monkey survey response webhook details.
	router.HandleFunc("/survey-response", routes.SurveyResponseWebhook).Methods(http.MethodPost)
	router.HandleFunc("/user/{id}/surveys", routes.GetUserSurveys).Methods(http.MethodGet)
	router.HandleFunc("/user/{id}/survey/{survey_id}/details", routes.GetUserSurveyDetails).Methods(http.MethodGet)
	router.HandleFunc("/save-token", routes.SaveSurveyMonkeyAccessToken).Methods(http.MethodPost)

	// add middleware
	authMiddleware := server.ValidateAccessToken()
	router.Use(authMiddleware.Handler)

	svr, err := server.DefaultServer()
	if err != nil {
		log.Fatalln(err)
	}

	svr.SetRouter(router)

	svr.ListenAndServe()

}
