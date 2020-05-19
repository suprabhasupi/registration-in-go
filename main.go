package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/subosito/gotenv"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// These are the three models atruct which has been created
type User struct {
	ID       int    `json: "id:`
	Email    string `json: "email"`
	Password string `json: "password"`
}

type JST struct {
	Token string `json : "token"`
}

type Error struct {
	Msg string `json : "msg"`
}

type JWT struct {
	Token string `json:"token"`
}

// creating new variable db
var db *sql.DB

func init() {
	gotenv.Load()
}
func main() {
	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))

	if err != nil {
		log.Fatal(err)
	}

	// invoke the open sql
	// first parameter is driver name
	db, err = sql.Open("postgres", pgUrl)
	// fmt.Println("db--->", db) // we will get all details of DB
	if err != nil {
		log.Fatal(err)
	}

	// Ping will return once value, if there is a connection established successfully the response,will return nil
	err = db.Ping()

	router := mux.NewRouter()
	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/protected", TokenVerifyMiddleware(protectedEndpoint)).Methods("GET")
	// why? TokenVerifyMiddleware(protectedEndpoint) => because protectedEndpoint needs token which will be generated from TokenVerifyMiddleware using JWT

	log.Println("Server is running on 9000...")
	// 1st parmeter is address and 2nd parameter handler function
	log.Fatal(http.ListenAndServe(":9000", router))
	// if there will be any error while starting the servee
}

func respondWithError(w http.ResponseWriter, status int, error Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

func reponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func signup(w http.ResponseWriter, r *http.Request) {
	// creating user variable by User struct type
	var user User
	var error Error
	// Decode returns an error
	json.NewDecoder(r.Body).Decode(&user)

	// Will do validation here, so email and password should not be empty
	if user.Email == "" {
		error.Msg = "Email is missing!"
		// http.StatusBadRequest is for 400 which is bad request
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Msg = "Password is missing!"
		respondWithError(w, http.StatusBadRequest, error)
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
		respondWithError(w, http.StatusInternalServerError, error)
		return
	}
	// if its success then return password as empty password , bcoz we don't want to send password explicit even its hashed and it cannot be reversible
	user.Password = ""
	// set Header
	w.Header().Set("Content-Type", "application/json")
	reponseJSON(w, user)
}

func GenerateToken(user User) (string, error) {
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

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login called")
	// w.Write([]byte("successfully caled Login")) // will get this as a response

	var user User
	// JWT will holding the token
	var jwt JWT
	// error message will snet to client
	var error Error
	// decoded the response body
	json.NewDecoder(r.Body).Decode(&user)
	// validating the user amila nd password
	if user.Email == "" {
		error.Msg = "Email is missing!"
		// http.StatusBadRequest is for 400 which is bad request
		respondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Msg = "Password is missing!"
		respondWithError(w, http.StatusBadRequest, error)
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
			respondWithError(w, http.StatusBadRequest, error)
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
		respondWithError(w, http.StatusUnauthorized, error)
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)

	jwt.Token = token
	reponseJSON(w, jwt)

	// fmt.Println(token)
}
func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("protectedEndpoint called")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
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
				respondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
			// spew.Dump(token) // we will get valid true here
			if token.Valid {
				// inbvoke the function that we passed into middleware
				next.ServeHTTP(w, r)
			} else {
				errorObject.Msg = error.Error()
				respondWithError(w, http.StatusUnauthorized, errorObject)
				return
			}
		} else {
			errorObject.Msg = "Invalid Token!"
			respondWithError(w, http.StatusUnauthorized, errorObject)
			return
		}
	})
}
