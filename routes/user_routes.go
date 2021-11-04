package routes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/apesurvey/ape-survey-backend/v2/models"
	"github.com/apesurvey/ape-survey-backend/v2/service"
	"github.com/apesurvey/ape-survey-backend/v2/utils"
	"github.com/gorilla/mux"
)

// DeleteUserHandler removes a user account from record.
// scopes: delete:users
func DeleteUserHandler(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	userID := params["id"]

	// TODO ensure the user does not have any active reward pools or rewards to claim.

	authService, err := service.NewAuthService()
	if err != nil {
		log.Println("error while creating new auth service instance: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	statusCode, err := authService.DeleteUser(userID)
	if err != nil {
		log.Println(": ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
}

func PatchUserHandler(w http.ResponseWriter, req *http.Request) {

	params := mux.Vars(req)
	userID := params["id"]

	// read request body
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(": ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	var requestBody models.PatchUserRequest

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Println(": ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	authService, err := service.NewAuthService()
	if err != nil {
		log.Println("error while creating new auth service instance: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	statusCode, err := authService.PatchUser(userID, requestBody.Email)
	if err != nil {
		log.Println(": ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	w.WriteHeader(statusCode)
}
