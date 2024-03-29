package main

import (
	"log"
	"net/http"
	"os"

	"github.com/apesurvey/ape-survey-backend/v2/routes"
	"github.com/apesurvey/ape-survey-backend/v2/server"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// routes
	router.HandleFunc("/user/{id}", routes.DeleteUserHandler).Methods(http.MethodDelete, http.MethodOptions)
	router.HandleFunc("/user/{id}", routes.PatchUserHandler).Methods(http.MethodPatch, http.MethodOptions)

	// TODO look up survey monkey survey response webhook details.
	router.HandleFunc("/survey-response", routes.SurveyResponseWebhook).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/user/{id}/surveys", routes.GetUserSurveys).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/user/{id}/survey/{survey_id}/details", routes.GetUserSurveyDetails).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/user/{id}/sm-connected", routes.SurveyMonkeyConnectionCheckHandler).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc("/oauth/token", routes.SurveyMonkeyOAuthToken).Methods(http.MethodPost, http.MethodOptions)

	// add middleware
	middleware := []mux.MiddlewareFunc{}
	authMiddleware := server.ValidateAccessToken()
	middleware = append(middleware, server.EnableCORS)

	if os.Getenv("DEV") != "true" {
		middleware = append(middleware, authMiddleware.Handler)
	}

	router.Use(middleware...)

	svr, err := server.DefaultServer()
	if err != nil {
		log.Fatalln(err)
	}

	svr.SetRouter(router)

	svr.ListenAndServe()

}
