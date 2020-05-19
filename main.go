package main

import (
	"database/sql"
	"log"
	"net/http"
	"registration-in-go/controllers"
	"registration-in-go/driver"

	"github.com/subosito/gotenv"

	"github.com/gorilla/mux"
)

// creating new variable db
var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	db = driver.ConnectDB()
	controllers := controllers.Controller{}
	router := mux.NewRouter()
	router.HandleFunc("/signup", controllers.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controllers.Login(db)).Methods("POST")
	router.HandleFunc("/protected", controllers.TokenVerifyMiddleware(controllers.protectedEndpoint())).Methods("GET")
	// why? TokenVerifyMiddleware(protectedEndpoint) => because protectedEndpoint needs token which will be generated from TokenVerifyMiddleware using JWT

	log.Println("Server is running on 9000...")
	// 1st parmeter is address and 2nd parameter handler function
	log.Fatal(http.ListenAndServe(":9000", router))
	// if there will be any error while starting the servee
}
