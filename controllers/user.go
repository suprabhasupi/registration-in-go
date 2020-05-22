package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"registration-in-go/models"
	"registration-in-go/utils"

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

		// userRepo := userRepository.UserRepository{}
		// user = userRepo.Signup(db, user)

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

// Login function
func (c Controller) Login(db *sql.DB) http.HandlerFunc {
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

		// userRepo := userRepository.UserRepository{}
		// user, err := userRepo.Login(db, user)

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
