package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kim3z/golang-rest-auth/mail"
	"github.com/kim3z/golang-rest-auth/models"
)

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

// ForgotPassword ...
func ForgotPassword(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	email := ps.ByName("email")
	w.Header().Add("Content-Type", "application/json")

	resetPswStatus, ok := models.SetResetPasswordToken(email)

	if !ok {
		json.NewEncoder(w).Encode(resetPswStatus)
		return
	}

	user := resetPswStatus["user"].(*models.User)
	fmt.Println(user.Email)

	resetLink := fmt.Sprintf("http://localhost:8080/api/user/reset/%s", user.TokenReset)
	mailData := map[string]string{"link": resetLink}

	mail.Send([]string{user.Email}, "Password reset", mailData, "resetpassword.html")
}
