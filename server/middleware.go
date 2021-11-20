package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/joho/godotenv"
)

var audience string
var issuer string
var clientURL string

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Panicln("unable to load .env file: ", err)
	}

	var ok bool
	audience, ok = os.LookupEnv("AUTH0_API")
	if !ok {
		log.Panicln("unable to get api identifier")
	}

	issuer, ok = os.LookupEnv("AUTH0_APP_DOMAIN")
	if !ok {
		log.Panicln("unable to get auth0 app domain")
	}

	clientURL, ok = os.LookupEnv("CLIENT_URL")
	if !ok {
		log.Panicln("unable to get client url")
	}

}

// Jwks holds the slice of JSON Web Keys.
type Jwks struct {
	Keys []JSONWebKey `json:"keys"`
}

// JSONWebKey represents the JWK object used to verify JWT's provided in the Authorization header of requests to the gateway.
type JSONWebKey struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// ValidateAccessToken examines the passed in bearer token in the authorization header of the request
// and ensures that it is valid and that the user has the required scopes to access the endpoint
func ValidateAccessToken() *jwtmiddleware.JWTMiddleware {
	return jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {

			// Verify 'aud' claim
			for k, v := range token.Claims.(jwt.MapClaims) {
				if k == "aud" {
					checkAud := false
					for _, aud := range v.([]interface{}) {
						if aud.(string) == audience+"/" {
							checkAud = true
							break
						}
					}

					if !checkAud {
						log.Printf("Invalid aud request attempt: %v\n", token.Claims)
						return token, errors.New("not authorized")
					}
				}
			}

			// Verify 'iss' claim
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(issuer, false)
			if !checkIss {
				log.Printf("Invalid iss request attempt: %v", token.Claims)
				return token, errors.New("not authorized")
			}

			cert, err := getPemCert(token)
			if err != nil {
				return nil, err
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

			return result, nil

		},
		SigningMethod: jwt.SigningMethodRS256,
	})
}

// getPemCert retrieves the remote JWKS for our Auth0 account and returns the certificate with the public key in PEM format.
func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	client := http.DefaultClient

	resp, err := client.Get(issuer + ".well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil

}

// EnableCORS allows the setting of the access control origin to restrict the sources
// that can send requests to the server.
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", clientURL)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		next.ServeHTTP(w, req)
	})
}
