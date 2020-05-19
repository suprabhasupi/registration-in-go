package controller

import (
	"fmt"
	"net/http"
	"registration-in-go/utils"
)

type Controller struct {
}

// ProtectedEndpoint function
func (c Controller) ProtectedEndpoint(w http.ResponseWriter, r *http.Request) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("protectedEndpoint called")
		utils.ResponseJSON(w, "yes")
	}
}
