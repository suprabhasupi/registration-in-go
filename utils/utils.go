package utils

import (
	"encoding/json"
	"net/http"
	"registration-in-go/models"
)

// RespondWithError for error case
func RespondWithError(w http.ResponseWriter, status int, message string) {
	var error models.Error
	error.Msg = message
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

// ResponseJSON in JSO format
func ResponseJSON(w http.ResponseWriter, data interface{}) {
	// set Header
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
