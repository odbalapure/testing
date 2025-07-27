package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// apply logic before hitting any route
	// Eg: when you get a request and something breaks rather than killing the entire application
	// recover from that and send an error message to the end user
	mux.Use(middleware.Recoverer)
	mux.Use(app.addIPToContext)
	mux.Use(app.Session.LoadAndSave)

	// register routes
	// mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })

	mux.Get("/", app.Home)
	mux.Post("/login", app.Login)
	mux.Get("/user/profile", app.Profile)

	// server static assets
	// `FileServer` returns a `Handler`
	// It reads files from the disk and return it over HTTP with headers
	// NOTE: If a request comes in for /static/logo.png, fileServer will look for: ./static/static/logo.png
	// Hence we need to strip `/static`
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
