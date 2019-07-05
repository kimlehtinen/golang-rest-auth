package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kim3z/golang-rest-auth/models"
)

// Authenticate ...
func Authenticate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Auth!\n")
}

// CreateUser ...
func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &models.User{}

	w.Header().Add("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := map[string]interface{}{
			"status":  false,
			"message": "User could not be created",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewEncoder(w).Encode(user.Create())
}

// LoginUser ...
func LoginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := &models.User{}

	w.Header().Add("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		resp := map[string]interface{}{
			"status":  false,
			"message": "Login failed",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Login(user.Email, user.Password)
	json.NewEncoder(w).Encode(resp)
}
