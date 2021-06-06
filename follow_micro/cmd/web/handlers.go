package main

import (
	"armani_follow/pkg/models"
	notify "armani_follow/protos"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strings"
	"time"
)



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



func (app *application) getFollowers(w http.ResponseWriter, r *http.Request)  {
	claims, _ := extractClaims(r)
	name, _  :=  claims["name"].(string)
	email, _ := claims["email"].(string)
	id, _ := claims["user_id"].(float64)
	me := models.User{
		ID:    id,
		Name:  name,
		Email: email,
	}

	stuf, err := app.followers.GetFollowers(me)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}else{
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(stuf)
	}
}


func (app *application) followUser(w http.ResponseWriter, r *http.Request)  {
	var person models.User
	var error models.Error

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		error.Message = err.Error()
		http.Error(w, error.Message,  http.StatusBadRequest)
		return
	}

	claims, _ := extractClaims(r)
	name, _  :=  claims["name"].(string)
	email, _ := claims["email"].(string)
	id, _ := claims["user_id"].(float64)

	if id == person.ID{
		error.Message = "You can not follow your self!!!"
		http.Error(w, error.Message,  http.StatusBadRequest)
		return
	}else{
		me := models.User{
			ID:    id,
			Name:  name,
			Email: email,
		}
		err = app.followers.Follow(me, person)

		if err != nil{
			error.Message = err.Error()
			http.Error(w, error.Message,  http.StatusBadRequest)
		}else{
			sendEmail(me.Email, person.Email)
			// check
			// rewrite
		}

	}
}

func (app *application) unfollowUser(w http.ResponseWriter, r *http.Request)  {
	var person models.User
	var error models.Error

	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		error.Message = err.Error()
		http.Error(w, error.Message,  http.StatusBadRequest)
		return
	}

	claims, _ := extractClaims(r)
	name, _  :=  claims["name"].(string)
	email, _ := claims["email"].(string)
	id, _ := claims["user_id"].(float64)

	if id == person.ID{
		error.Message = "You can not Unfollow your self!!!"
		http.Error(w, error.Message,  http.StatusBadRequest)
		return
	}else{
		me := models.User{
			ID:    id,
			Name:  name,
			Email: email,
		}
		err = app.followers.Unfollow(me, person)

		if err != nil{
			error.Message = err.Error()
			http.Error(w, error.Message,  http.StatusBadRequest)
		}else{
			sendEmail(me.Email, person.Email)
			// check
			// rewrite
		}

	}


}

func sendEmail(from,to string)  {



	requests := notify.NotificationRequest{
		From: from,
		To:   to,
	}


	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Println("could not connect: %v", err)
	}
	defer conn.Close()

	c := notify.NewNotificationServiceClient(conn)


	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err = c.NotifyFollow(ctx, &requests)
	if err != nil {
		log.Println("error while calling LongGreet: %v", err)
	}

}