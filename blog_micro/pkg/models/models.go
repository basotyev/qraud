package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// Add a new ErrInvalidCredentials error. We'll use this later if a user
	//tries to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// Add a new ErrDuplicateEmail error. We'll use this later if a user
	// tries to signup with an email address that's already in use.
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

type Snippet struct {
	ID int					`json:"id"`
	Title string			`json:"title"`
	Content string			`json:"content"`
	Created time.Time		`json:"created"`
	UserId int 				`json:"user_id"`
	UserName string			`json:"user_name"`
}


type Error struct {
	Message string `json:"message"`
}