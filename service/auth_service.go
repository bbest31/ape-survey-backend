package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/apesurvey/ape-survey-backend/v2/constants"
	"github.com/apesurvey/ape-survey-backend/v2/models"
	"github.com/apesurvey/ape-survey-backend/v2/utils"
	"github.com/joho/godotenv"
)

type AuthService struct {
	Client  *http.Client
	AuthURL string
	Token   string
}

func NewAuthService() (AuthService, error) {

	err := godotenv.Load()
	if err != nil {
		return AuthService{}, fmt.Errorf("unable to load .env file")
	}

	token, ok := os.LookupEnv("AUTH0_MANAGEMENT_API_TOKEN")
	if !ok || utils.IsEmptyString(token) {
		return AuthService{}, fmt.Errorf("authentication api url not sent in environment")
	}

	client := http.DefaultClient

	return AuthService{client, constants.AUTH0_API, token}, nil

}

func (service AuthService) DeleteUser(userID string) (int, error) {

	r, err := http.NewRequest("DELETE", service.AuthURL+"/users/"+userID, nil)
	if err != nil {
		return 0, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", service.Token))

	resp, err := service.Client.Do(r)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

// PatchUser sends the request to update information of a user.
// Currently we only support the updating of their email.
// scopes: update:users update:users_app_metadata
func (service AuthService) PatchUser(userID string, email string) (int, error) {

	body, err := json.Marshal(models.PatchUserRequest{Email: email, ClientID: constants.AUTH0_CLIENT_ID})
	if err != nil {
		return 0, err
	}

	requestBody := bytes.NewBuffer(body)
	r, err := http.NewRequest("PATCH", service.AuthURL+"/users/"+userID, requestBody)
	if err != nil {
		return 0, err
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %v", service.Token))

	resp, err := service.Client.Do(r)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}

func (service AuthService) Close() {
	service.Client.CloseIdleConnections()
}
