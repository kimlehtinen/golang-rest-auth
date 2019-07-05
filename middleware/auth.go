package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/kim3z/golang-rest-auth/models"
)

var respMissingToken = map[string]interface{}{
	"status":  false,
	"message": "Missing token",
}

var respInvalidToken = map[string]interface{}{
	"status":  false,
	"message": "Invalid token",
}

var respMalformedToken = map[string]interface{}{
	"status":  false,
	"message": "Malformed token",
}

// Authenticate request
func Authenticate(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		headerToken := r.Header.Get("Authorization") // get jwt token

		if headerToken == "" {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respMissingToken)
			return
		}

		headerTokenData := strings.Split(headerToken, " ") // get token from bearer
		if len(headerTokenData) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respInvalidToken)
			return
		}

		htdToken := headerTokenData[1]
		jwtT := &models.JwtToken{}
		token, err := jwt.ParseWithClaims(htdToken, jwtT, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("jwt_secret")), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respMalformedToken)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(respInvalidToken)
			return
		}

		ctx := context.WithValue(r.Context(), "user", jwtT.UserID)
		r = r.WithContext(ctx)

		next(w, r, ps)
	})
}
