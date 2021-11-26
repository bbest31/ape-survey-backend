package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/apesurvey/ape-survey-backend/v2/constants"
	"github.com/apesurvey/ape-survey-backend/v2/models"
	"github.com/apesurvey/ape-survey-backend/v2/service"
	"github.com/apesurvey/ape-survey-backend/v2/utils"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"google.golang.org/grpc/codes"
)

// SaveSurveyMonkeyAccessToken persists the use access token for using the SurveyMonkey API into GCP Secret Manager.
// https://cloud.google.com/secret-manager/docs/creating-and-accessing-secrets#add-secret-version
func SaveSurveyMonkeyAccessToken(w http.ResponseWriter, req *http.Request) {

	// validate request body details
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("error while reading save-token req body: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	var requestBody models.SaveTokenRequest
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Println(" /save-token unable to unmarhsal request body: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	// save the token in the Secret Manager
	ctx := context.Background()
	secretManagerService, err := service.NewClient(ctx)
	if err != nil {
		log.Println("error while building Secret Manager client: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	defer secretManagerService.Close()

	err = secretManagerService.CreateSecretRequest(ctx, requestBody.UserID, requestBody.AccessToken)
	if err != nil {
		log.Println("error while saving secret: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	w.WriteHeader(200)

}

// TODO - implement
func SurveyResponseWebhook(w http.ResponseWriter, req *http.Request) {

}

// GetUserSurveys returns a list of surveys owned or shared with the authenticated user.
// This SurveyMonkey endpoint needs the View Surveys scope.
func GetUserSurveys(w http.ResponseWriter, req *http.Request) {

	// get user access token if connected to SurveyMonkey account
	ctx := context.Background()
	secretManagerService, err := service.NewClient(ctx)
	if err != nil {
		log.Println("error while building Secret Manager client: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	params := mux.Vars(req)
	userID := params["id"]

	secret, err := secretManagerService.AccessSecret(fmt.Sprintf("projects/%s/secrets/%s/versions/latest", constants.GCP_PROJECT_ID, userID), ctx)
	if err != nil {
		log.Println("error while requesting user SM access token: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}
	accessToken := string(secret)

	// request for user surveys
	client := &http.Client{}
	defer client.CloseIdleConnections()

	r, err := http.NewRequest("GET", constants.SURVEY_MONKEY_API+"/surveys", nil)
	if err != nil {
		log.Println("error while creating new request for user surveys: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	resp, err := client.Do(r)
	if err != nil {
		log.Println("error while requesting user surveys: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	surveys := []models.Survey{}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error while reading SM repsonse: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	// read the total attribute to see how many more pages there are.
	data := models.SurveysResponse{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("error while unmarshalling response: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	surveys = append(surveys, data.Data...)

	// for every unread page of surveys send a request and add data to the result
	if pages := data.Total; pages > 1 {
		for i := 2; i <= pages+1; i++ {
			r, err = http.NewRequest("GET", constants.SURVEY_MONKEY_API+fmt.Sprintf("/surveys?page=%v", i), nil)
			if err != nil {
				log.Println("error while creating new request for user surveys, leaving some out: ", err)
			}

			r.Header.Add("Accept", "application/json")
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

			res, err := client.Do(r)
			if err != nil {
				log.Println("error while requesting user surveys: ", err)
			}

			body, err = ioutil.ReadAll(res.Body)
			if err != nil {
				log.Println("error while reading SM repsonse: ", err)
				utils.SendErrorResponse(w, err.Error())
				return
			}

			data := models.SurveysResponse{}

			err = json.Unmarshal(body, &data)
			if err != nil {
				log.Println("error while unmarshalling response: ", err)
				utils.SendErrorResponse(w, err.Error())
				return
			}

			surveys = append(surveys, data.Data...)

		}
	}

	// return the list of surveys to the frontend.

	err = utils.SendResponseWithData(w, surveys)
	if err != nil {
		log.Println("error while packaging response: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

}

// GetUserSurveyDetails retrieves the question bank for a specific survey.
func GetUserSurveyDetails(w http.ResponseWriter, req *http.Request) {

	path := strings.Split(req.URL.Path, "/")
	userID := path[2]
	surveyID := path[4]

	fmt.Printf("User id = %v and Survey ID = %v", userID, surveyID)

	client := &http.Client{}
	defer client.CloseIdleConnections()

	r, err := http.NewRequest("GET", constants.SURVEY_MONKEY_API+fmt.Sprintf("/surveys/%v/details", surveyID), nil)
	if err != nil {
		log.Println("error while creating new request for user surveys: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	// get user access token if connected to SurveyMonkey account
	ctx := context.Background()
	secretManagerService, err := service.NewClient(ctx)
	if err != nil {
		log.Println("error while building Secret Manager client: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	secret, err := secretManagerService.AccessSecret(fmt.Sprintf("projects/%s/secrets/%s/versions/latest", constants.GCP_PROJECT_ID, userID), ctx)
	if err != nil {
		log.Println("error while requesting user SM access token: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}
	accessToken := string(secret)

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	// send request
	res, err := client.Do(r)
	if err != nil {
		log.Println("error while requesting user surveys: ", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error while reading SM repsonse: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	//package response
	data := models.SurveyDetailsResponse{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("error while unmarshalling response: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	questions := []string{}

	for _, page := range data.Pages {
		questions = append(questions, page.Questions...)
	}

	err = utils.SendResponseWithData(w, questions)
	if err != nil {
		log.Println("error while packaging response: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}
}

// SurveyMonkeyConnectionCheckHandler checks if the user has connected their SurveyMonkey account.
func SurveyMonkeyConnectionCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	userID := params["id"]

	// get user access token if connected to SurveyMonkey account
	ctx := context.Background()
	secretManagerService, err := service.NewClient(ctx)
	if err != nil {
		log.Println("error while building Secret Manager client: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	_, err = secretManagerService.AccessSecret(fmt.Sprintf("projects/%s/secrets/%s/versions/latest", constants.GCP_PROJECT_ID, userID), ctx)
	if err != nil {
		if err.(*apierror.APIError).GRPCStatus().Code() == codes.NotFound {
			// user has not connected SM account
			response := struct {
				SMConnected bool `json:"sm_connected`
			}{
				SMConnected: false,
			}

			json.NewEncoder(w).Encode(response)
		} else {
			log.Println("error while requesting user SM access token: ", err)
			utils.SendErrorResponse(w, err.Error())
		}
		return
	}

	response := struct {
		SMConnected bool `json:"sm_connected`
	}{
		SMConnected: true,
	}

	json.NewEncoder(w).Encode(response)

}

func SurveyMonkeyOAuthToken(w http.ResponseWriter, req *http.Request) {

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("error while reading oauth/token request body: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	var tokenRequest models.OAuthTokenRequest

	err = json.Unmarshal(data, &tokenRequest)
	if err != nil {
		log.Println("error while unmarshalling request body: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	client := &http.Client{}
	defer client.CloseIdleConnections()

	err = godotenv.Load()
	if err != nil {
		log.Println("error while reading .env file: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	clientID, ok := os.LookupEnv("SM_CLIENT_ID")
	if !ok {
		log.Println("error while reading client id environment variable: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	clientSecret, ok := os.LookupEnv("SM_CLIENT_SECRET")
	if !ok {
		log.Println("error while reading client secret environment variable: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	requestData := url.Values{}
	requestData.Set("redirect_uri", "http://localhost:3000/dashboard/oauth2callback") // TODO make this progammatic
	requestData.Set("client_secret", clientSecret)
	requestData.Set("client_id", clientID)
	requestData.Set("code", tokenRequest.Code)
	requestData.Set("grant_type", "authorization_code")

	r, err := http.NewRequest("POST", "https://api.surveymonkey.com/oauth/token", strings.NewReader(requestData.Encode()))
	if err != nil {
		log.Println("error while creating new request for oath token: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(requestData.Encode())))

	// send request
	res, err := client.Do(r)
	if err != nil {
		log.Println("error while requesting oauth token: ", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error while reading SM response: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	var tokenResponse models.OAuthTokenResponse

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Println("error while unmarshalling response body: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	// Save access token in Secret Manager.
	ctx := context.Background()
	secretManagerService, err := service.NewClient(ctx)
	if err != nil {
		log.Println("error while building Secret Manager client: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	defer secretManagerService.Close()

	err = secretManagerService.CreateSecretRequest(ctx, tokenRequest.UserID, tokenResponse.AccessToken)
	if err != nil {
		log.Println("error while saving secret: ", err)
		utils.SendErrorResponse(w, err.Error())
		return
	}

	w.WriteHeader(200)

}
