package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"registration-in-go/models"

	"github.com/dgrijalva/jwt-go"
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

func GenerateToken(user models.User) (string, error) {
	var err error
	// assign any variable
	secret := os.Getenv("SECRET")
	// there are many claims, which we can go through the docs of JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"iss":   "course",
	})

	// spew.Dump(token)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}
