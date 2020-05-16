package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
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
	msg string
}

// creating new variable db
var db *sql.DB

func main() {
	pgUrl, err := pq.ParseURL("POSTGRESURL")

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

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signup called")
	w.Write([]byte("successfully caled Signup")) // will get this as a response
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("login called")
	w.Write([]byte("successfully caled Login")) // will get this as a response
}
func protectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("TokenVerifyMiddleware called")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	fmt.Println("TokenVerifyMiddleware called")
	return nil
}
