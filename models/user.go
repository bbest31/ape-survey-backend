package models

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
}
