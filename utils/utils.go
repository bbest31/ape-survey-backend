package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
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

func IsEmptyString(s string) bool {
	return strings.Trim(s, " ") == ""
}

// InterfaceSlice makes a slice of any type into a slice of interface
func InterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, fmt.Errorf("parameter is a non-slice type")
	}

	if s.IsNil() {
		return nil, nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}
