package main

import (
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	Session *scs.SessionManager
}

func main() {
	// setup an application config
	app := application{}

	// get a session manager
	// NOTE: setup the session before hitting the routes
	app.Session = getSession()

	// setup application routes
	mux := app.routes()

	// print out a message
	log.Println("Strating serve on port 8080")

	// start the server
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
