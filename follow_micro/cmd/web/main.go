package main

import (

	"armani_follow/pkg/models/myMongo"
	"context"
	"flag"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"

	"log"
	"net/http"
)

type application struct {
	errorLog 			*log.Logger
	infoLog 			*log.Logger
	followers 			*myMongo.Followers
}



func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Method("POST", "/follow", TokenVerifyMiddleWare(app.followUser))
	mux.Method("POST", "/unfollow", TokenVerifyMiddleWare(app.unfollowUser))
	mux.Method("POST", "/getfollower", TokenVerifyMiddleWare(app.getFollowers))
	return mux
}



func main() {

	addr := flag.String("addr", ":4002", "HTTP network address")

	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client1, _ := mongo.Connect(ctx, clientOptions)




	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		followers: &myMongo.Followers{client1},
	}




	srv := &http.Server{
	Addr: *addr,
	ErrorLog: errorLog,
	Handler: app.routes(),
	IdleTimeout: time.Minute,
	ReadTimeout: 5 * time.Second,
	WriteTimeout: 10 * time.Second,
	}



	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}