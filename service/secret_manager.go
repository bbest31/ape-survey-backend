package service

import (
	"context"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/apesurvey/ape-survey-backend/v2/constants"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManagerService struct {
	Client *secretmanager.Client
}

func NewClient(ctx context.Context) (SecretManagerService, error) {

	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return SecretManagerService{}, err
	}

	service := SecretManagerService{client}

	return service, nil
}

func (service SecretManagerService) CreateSecretRequest(ctx context.Context, secretID string, data string) error {
	createSecretReq := secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", constants.GCP_PROJECT_ID),
		SecretId: secretID,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
		},
	}

	secret, err := service.Client.CreateSecret(ctx, &createSecretReq)
	if err != nil {
		return err
	}

	// Declare the payload to store
	payload := []byte(data)

	// Build the request
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: secret.Name,
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	// Call the API
	_, err = service.Client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return err
	}

	return nil
}

func (service SecretManagerService) AccessSecret(versionName string, ctx context.Context) ([]byte, error) {
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: versionName,
	}

	result, err := service.Client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, err
	}

	return result.Payload.GetData(), nil
}

func (service SecretManagerService) Close() {
	service.Client.Close()
}
