package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrInvalidCredentials = errors.New("models: invalid credentials")

	ErrDuplicateEmail = errors.New("models: duplicate email")
)



type User struct {
	ID 		int							`json:"id"`
	Name 	string						`json:"name"`
	Email 	string						`json:"email"`
	HashedPassword string		    	`json:"hashed_password"`
	Created time.Time					`json:"created"`
	Active bool							`json:"active"`
}

type JWT struct {
	Token string			`json:"token"`
}

type Error struct {
	Message string 			`json:"message"`
}


