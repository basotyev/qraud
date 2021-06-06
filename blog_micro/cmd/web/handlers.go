package main

import (
	"armani_blog/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
	"strings"
)


func ResponseJSON(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}


func (app *application) showSnippetLatest(w http.ResponseWriter, r *http.Request)  {
	var error models.Error
	s, err := app.snippets.Latest()
	if err != nil {
		error.Message = err.Error()
		RespondWithError(w, http.StatusBadRequest, error)
	}
	ResponseJSON(w, s)
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request)  {
	id := chi.URLParam(r, "id")
	intid, err := strconv.Atoi(id)
	if id == "" || err != nil || intid < 1{
		app.errorLog.Println("id not supported")
		return
	}
	s, err := app.snippets.GetByUserId(intid)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.errorLog.Println(models.ErrNoRecord.Error())
		} else {
			app.errorLog.Println(err)
		}
		return
	}
	ResponseJSON(w, s)
	fmt.Println(s)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request)  {
	claims, _ := extractClaims(r)
	name, _  :=  claims["name"].(string)
	//email, _ := claims["email"].(string)
	id, _ := claims["user_id"].(float64)
	var snip models.Snippet
	json.NewDecoder(r.Body).Decode(&snip)

	_ , err := app.snippets.Insert(snip.Title, snip.Content, name, int(id) )

	if err != nil{
		app.errorLog.Println(err.Error())
	}
}

func extractClaims(r *http.Request) (jwt.MapClaims, bool) {
	authHeader := r.Header.Get("Authorization")
	bearerToken := strings.Split(authHeader, " ")
	tokenStr := bearerToken[1]

	hmacSecretString := "astanaIT"
	hmacSecret := []byte(hmacSecretString)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return hmacSecret, nil
	})

	if err != nil {
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}


func (app *application) updateSnippet(w http.ResponseWriter, r *http.Request)  {
	claims, _ := extractClaims(r)
	//name, _  :=  claims["name"].(string)
	id, _ := claims["user_id"].(float64)
	var snip models.Snippet

	err:= json.NewDecoder(r.Body).Decode(&snip)
	if err != nil{
		app.errorLog.Println(err.Error())
		return
	}

	s, err := app.snippets.UpdateSnippet(snip.ID, int(id) , &snip)

	if err != nil{
		app.errorLog.Println(err.Error())
	}else {
		ResponseJSON(w, s)
	}


}