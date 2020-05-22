package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"registration-in-go/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// RespondWithError for error case
func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	// var error models.Error
	// error.Msg = message
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

// ResponseJSON in JSO format
func ResponseJSON(w http.ResponseWriter, data interface{}) {
	// set Header
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// GenerateToken function
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

// TokenVerifyMiddleWare function
func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyMiddleware called")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return to client any error that ecounter next
		var errorObject models.Error
		// holding the value of authorization header taht we send fromclient to server the request obj has a field called header
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			// to check the token is valid
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				// the value will be returned fropm oarse will be token
				return []byte(os.Getenv("SECRET")), nil
			})

			if error != nil {
				errorObject.Msg = error.Error()
				RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
			// spew.Dump(token) // we will get valid true here
			if token.Valid {
				// inbvoke the function that we passed into middleware
				next.ServeHTTP(w, r)
			} else {
				errorObject.Msg = error.Error()
				RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Msg = "Invalid Token!"
			RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
