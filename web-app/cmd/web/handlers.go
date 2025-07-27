package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"
	"webapp/pkg/data"
)

var pathToTempaltes = "./templates/"

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]any)

	// fmt.Fprint(w, "Home Page")
	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.Get(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "Hit this page "+time.Now().UTC().String())
	}

	err := app.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
	fmt.Print(err)
}

func (app *application) Profile(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "profile.page.gohtml", &TemplateData{})
}

type TemplateData struct {
	IP string
	// `any` was part of go v1.19
	// this is an alias for `interface`
	Data  map[string]any
	Error string
	Flash string
	User  data.User
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTempaltes, t), path.Join(pathToTempaltes+"base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	td.IP = app.ipFromContext(r.Context())
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")

	// execute the template, passing it data, if any
	err = parsedTemplate.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// validate data
	form := NewForm(r.PostForm)
	form.Required("email", "password")

	if !form.Valid() {
		app.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !form.Valid() {
		fmt.Fprint(w, "failed validation")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// get the user
	user, err := app.DB.GetUserByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid login!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	fmt.Println(email, password, user)

	// authenticate the user
	// redirect with error

	// prevent session fixation attack
	_ = app.Session.RenewToken(r.Context())

	// store success message in the session

	// redirect to some other page
	app.Session.Put(r.Context(), "flash", "Successfully logged in!")
	http.Redirect(w, r, "/user/profile", http.StatusSeeOther)

	fmt.Fprint(w, email)
}
