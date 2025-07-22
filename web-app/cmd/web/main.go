package main

import (
	"log"
	"net/http"
)

type application struct{}

func main() {
	// setup an application config
	app := application{}

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
