package models

type User struct {
	ID 		float64			`json:"id" bson:"id"`
	Name 	string			`json:"name" bson:"name"`
	Email 	string			`json:"email" bson:"email"`
}


type MongoResponse struct {
	IDUser 		float64					`json:"id" bson:"_id"`
	Subscribers []User 					`json:"subscribers" bson:"subscribers"`
	Subscribing []User					`json:"subscribing" bson:"subscribing"`
}

type JWT struct {
	Token string			`json:"token"`
}

type Error struct {
	Message string 			`json:"message"`
}


