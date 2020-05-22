package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"registration-in-go/utils"
)

// Controller model
type Controller struct{}

// ProtectedEndpoint function
func (c Controller) ProtectedEndpoint(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("protectedEndpoint called")
		utils.ResponseJSON(w, "yes")
	}
}
