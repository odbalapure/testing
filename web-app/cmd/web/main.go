package main

import (
	"flag"
	"log"
	"net/http"
	"webapp/pkg/db"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
	DSN     string
	DB      db.PostgresConn
}

func main() {
	// setup an application config
	app := application{}

	// when someone accesses `app.DSN` they get this though out their project
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	// When the main function exits, close the DB connection
	defer conn.Close()

	// app.DB = conn
	app.DB = db.PostgresConn{DB: conn}

	// get a session manager
	// NOTE: setup the session before hitting the routes
	app.Session = getSession()

	// setup application routes
	mux := app.routes()

	// print out a message
	log.Println("Strating serve on port 8080")

	// start the server
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
