package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
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

	resetLink := fmt.Sprintf("http://localhost:8080/api/user/reset-psw-check/%s", user.TokenReset)
	mailData := map[string]string{"link": resetLink}

	mail.Send([]string{user.Email}, "Password reset", mailData, "resetpassword.html")
}

// ResetPasswordCheck checks if user can change password
func ResetPasswordCheck(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	resetToken := ps.ByName("reset-token")

	// check if user with such token exists
	user := &models.User{}
	err := models.DB().Table("users").Where("token_reset = ? AND token_reset_expires > ?", resetToken, time.Now()).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			resp := map[string]interface{}{
				"status":  false,
				"message": "Either your reset token doesn't exist or it has expired",
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		resp := map[string]interface{}{
			"status":  false,
			"message": "DB error",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// check if token is valid
	jwtT := &models.JwtToken{}
	token, err := jwt.ParseWithClaims(resetToken, jwtT, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("jwt_secret")), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		resp := map[string]interface{}{
			"status":  false,
			"message": "Malformed token",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusForbidden)
		resp := map[string]interface{}{
			"status":  false,
			"message": "Invalid token",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	user.Password = ""
	resp := map[string]interface{}{
		"status":  true,
		"message": "Token is valid, you can change your password!",
		"user":    user,
	}
	json.NewEncoder(w).Encode(resp)
}
