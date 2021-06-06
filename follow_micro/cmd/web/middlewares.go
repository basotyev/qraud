package main

import (
	"armani_follow/pkg/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
)

func TokenVerifyMiddleWare(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorObject models.Error
		authHeader := r.Header.Get("Authorization")
		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]

			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}

				return []byte("astanaIT"), nil
			})

			if error != nil {
				errorObject.Message = error.Error()
				http.Error(w,  errorObject.Message, http.StatusUnauthorized)
				return
			}

			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				errorObject.Message = error.Error()
				http.Error(w,  errorObject.Message, http.StatusUnauthorized)
				return
			}
		} else {
			errorObject.Message = "Invalid token."
			http.Error(w,  errorObject.Message, http.StatusUnauthorized)
			return
		}
	})
}



