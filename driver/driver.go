package driver

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

// creating new variable db
var db *sql.DB

// Coonecting DB
func ConnectDB() *sql.DB {
	pgURL, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))

	if err != nil {
		log.Fatal(err)
	}

	// invoke the open sql
	// first parameter is driver name
	db, err = sql.Open("postgres", pgURL)
	// fmt.Println("db--->", db) // we will get all details of DB
	if err != nil {
		log.Fatal(err)
	}

	// Ping will return once value, if there is a connection established successfully the response,will return nil
	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return db
}
