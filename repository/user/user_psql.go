package userRepository

import (
	"database/sql"
	"log"
	"registration-in-go/models"
)

// UserRepository model
type UserRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Signup query
func (u UserRepository) Signup(db *sql.DB, user models.User) models.User {
	statement := "insert into users (email,password) values ($1, $2) RETURNING id;"
	// QueryRow will execute one row, Scan is supposed to return an error or nil. If QueryRow will not select any row then it will throw error
	err := db.QueryRow(statement, user.Email, user.Password).Scan(&user.ID)

	logFatal(err)

	user.Password = ""
	return user
}

// Login query
func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {
	row := db.QueryRow("select * from users where email=$1", user.Email)
	err := row.Scan(&user.ID, &user.Email, &user.Password)

	if err != nil {
		return user, err
	}

	return user, nil
}

// the methods are going to be reflective of db tabel interactives code, that we have inside our user controller
