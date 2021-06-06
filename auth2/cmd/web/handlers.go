package main

import (
	"armani_auth/pkg/models"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
)



func GenerateToken(user models.User) (string, error) {
	var err error
	secret := "hard_work_pays_off"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"user_id" : float64(user.ID),
		"name" : user.Name,
	})

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}

	return tokenString, nil
}




func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	var jwt models.JWT
	var error models.Error
	//var id int


	json.NewDecoder(r.Body).Decode(&user)


	if user.Email == "" {
		error.Message = "Email is missing."
		RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	if user.HashedPassword == "" {
		error.Message = "Password is missing."
		RespondWithError(w, http.StatusBadRequest, error)
		return
	}

	user, err := app.users.Authenticate(user.Email, user.HashedPassword)


	if err != nil {
		error.Message = err.Error()
		RespondWithError(w, http.StatusBadRequest, error)

	}else{
		token, err := GenerateToken(user)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		jwt.Token = token

		ResponseJSON(w, jwt)
	}
}



func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {



	var user models.User
	var error1 models.Error
	json.NewDecoder(r.Body).Decode(&user)

	if user.Email == "" {
		error1.Message = "Email is missing."
		RespondWithError(w, http.StatusBadRequest, error1)
		return
	}

	if user.HashedPassword == "" {
		error1.Message = "Password is missing."
		RespondWithError(w, http.StatusBadRequest, error1)
		return
	}


	err := app.users.Insert(user.Name, user.Email, user.HashedPassword)
	if err != nil {
		error1.Message = err.Error()
		RespondWithError(w, http.StatusInternalServerError, error1)
		return
	}else{
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(true)
	}

}

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}


func RespondWithError(w http.ResponseWriter, status int, error models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(error)
}

