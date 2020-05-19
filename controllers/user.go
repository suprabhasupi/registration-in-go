package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"registration-in-go/models"
	"registration-in-go/utils"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Signup and login are expecting db, wo we amking function to get the db
func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// creating user variable by User struct type
		var user models.User
		var error models.Error
		// Decode returns an error
		json.NewDecoder(r.Body).Decode(&user)

		// Will do validation here, so email and password should not be empty
		if user.Email == "" {
			error.Msg = "Email is missing!"
			// http.StatusBadRequest is for 400 which is bad request
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Msg = "Password is missing!"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		// bytes which represent the user password and cost will be 10
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		// converting bytes into string
		user.Password = string(hash)

		statement := "insert into users (email,password) values ($1, $2) RETURNING id;"
		// QueryRow will execute one row, Scan is supposed to return an error or nil. If QueryRow will not select any row then it will throw error
		err = db.QueryRow(statement, user.Email, user.Password).Scan(&user.ID)

		if err != nil {
			error.Msg = "Server Error"
			utils.RespondWithError(w, http.StatusInternalServerError, error)
			return
		}
		// if its success then return password as empty password , bcoz we don't want to send password explicit even its hashed and it cannot be reversible
		user.Password = ""
		utils.ResponseJSON(w, user)
	}
}

// Login
func (c, Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login called")
		// w.Write([]byte("successfully caled Login")) // will get this as a response

		var user models.User
		// JWT will holding the token
		var jwt models.JWT
		// error message will snet to client
		var error models.Error
		// decoded the response body
		json.NewDecoder(r.Body).Decode(&user)
		// validating the user amila nd password
		if user.Email == "" {
			error.Msg = "Email is missing!"
			// http.StatusBadRequest is for 400 which is bad request
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Msg = "Password is missing!"
			utils.RespondWithError(w, http.StatusBadRequest, error)
			return
		}
		//
		password := user.Password
		// check the user exist in DB by email, as email will be unique (also we can use ID)
		row := db.QueryRow("select * from users where email=$1", user.Email)
		err := row.Scan(&user.ID, &user.Email, &user.Password)

		if err != nil {
			if err == sql.ErrConnDone {
				error.Msg = "The User does not exist!"
				utils.RespondWithError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		// spew.Dump(user) // here the password will be in hashed which need to get decrypted

		hashedPassword := user.Password
		// copare the hashed password using bcrypt
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			error.Msg = "Invalid password!"
			utils.RespondWithError(w, http.StatusUnauthorized, error)
			return
		}

		token, err := utils.GenerateToken(user)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)

		jwt.Token = token
		utils.ResponseJSON(w, jwt)

		// fmt.Println(token)
	}

}

// TokenVerifyMiddleware function
func (c Controller) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyMiddleware called")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// return to client any error that ecounter next
		var errorObject Error
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
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
			// spew.Dump(token) // we will get valid true here
			if token.Valid {
				// inbvoke the function that we passed into middleware
				next.ServeHTTP(w, r)
			} else {
				errorObject.Msg = error.Error()
				utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Msg = "Invalid Token!"
			utils.RespondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
