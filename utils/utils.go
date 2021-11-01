package utils

import (
	"encoding/json"
	"net/http"
)

// SendErrorResponse sends a response using the http.ResponseWriter with status 500 and an attached json object containing the error message.
func SendErrorResponse(w http.ResponseWriter, message string) {

	data := map[string]interface{}{
		"ok":      false,
		"message": message,
	}

	body, err := json.Marshal(data)
	if err != nil {
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)

}

func SendResponseWithData(w http.ResponseWriter, data interface{}) error {
	payload := map[string]interface{}{
		"ok":   true,
		"data": data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
	return nil
}
