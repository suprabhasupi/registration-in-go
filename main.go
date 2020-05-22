package main

import (
	"database/sql"
	"log"
	"net/http"
	"registration-in-go/controllers"
	"registration-in-go/driver"
	"registration-in-go/utils"

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
	router := mux.NewRouter()

	controller := controllers.Controller{}

	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")
	router.HandleFunc("/protected", utils.TokenVerifyMiddleWare(controller.ProtectedEndpoint(db))).Methods("GET")

	// why? TokenVerifyMiddleware(protectedEndpoint) => because protectedEndpoint needs token which will be generated from TokenVerifyMiddleware using JWT

	log.Println("Server is running on 9000...")
	// 1st parmeter is address and 2nd parameter handler function
	log.Fatal(http.ListenAndServe(":9000", router))
	// if there will be any error while starting the servee
}
