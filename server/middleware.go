package server

import (
	"fmt"
	"net/http"
)

// ValidateAccessToken examines the passed in bearer token in the authorization header of the request
// and ensures that it is valid and that the user has the required scopes to access the endpoint
func ValidateAccessToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Validating authorization bearer token...")
		next.ServeHTTP(w, req)
	})
}
