package main

import (
	postgresql "armani_blog/pkg/models/postgres"
	"context"
	"flag"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"os"
	"time"
)




type application struct {
	errorLog 			*log.Logger
	infoLog 			*log.Logger
	snippets 			*postgresql.SnippetModel
}



func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Method("POST", "/create", TokenVerifyMiddleWare(app.createSnippet))
	mux.Method("PUT", "/update", TokenVerifyMiddleWare(app.updateSnippet))
	mux.Method("GET", "/user/{id}", TokenVerifyMiddleWare(app.showSnippet))
	mux.Method("GET", "/latest", TokenVerifyMiddleWare(app.showSnippetLatest))

	return mux
}

func OpenDB(dsn string) (*pgxpool.Pool, error)  {
	dbpool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}



func main() {

	addr := flag.String("addr", ":4001", "HTTP network address")
	dsn := flag.String("DBDsn", "postgres://postgres:0000@localhost:5432/blogs",
		"Connection string to DB")


	flag.Parse()




	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)



	dbpool, err := OpenDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer dbpool.Close()


	app := &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &postgresql.SnippetModel{DB: dbpool},

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
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}