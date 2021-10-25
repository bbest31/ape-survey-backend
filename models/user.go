package models

import "time"

// User is the model struct that represents an ApeSurvey user.
type User struct {
	// ID unique id for the ApeSurvey user.
	ID string `json:"id"`
	// Email is the users email address.
	Email string `json:"email"`
	// SurveyMonkeyLinked indicates if the users SurveyMonkey account has been linked.
	SurveyMonkeyLinked bool `json:"sm_linked"`
	// SurveyMonkeyAccessToken is a long-lived token used to utilize the SurveyMonkey integration.
	SurveyMonkeyAccessToken string `json:"sm_access_token"`
	// StripeCustomerID is used to link the user to the customer object in Stripe which saves their CC info.
	StripeCustomerID string `json:"stripe_customer_id"`
	// EmailNotfications indicate if the user wants to get an email notification when they earn a reward.
	EmailNotifications bool `json:"email_notifications"`
	// ShadowAccount indicates if the user has signed up on ApeSurvey or not.
	// If false this person has earned crypto via a survey response, but has not registered to claim it.
	ShadowAccount bool      `json:"shadow_account"`
	DeletedAt     time.Time `json:"deleted_at"`
}
