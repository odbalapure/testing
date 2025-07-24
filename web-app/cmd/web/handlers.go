package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
	"time"
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

type TemplateData struct {
	IP string
	// `any` was part of go v1.19
	// this is an alias for `interface`
	Data map[string]any
}

func (app *application) render(w http.ResponseWriter, r *http.Request, t string, data *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTempaltes, t), path.Join(pathToTempaltes+"base.layout.gohtml"))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	data.IP = app.ipFromContext(r.Context())

	// execute the template, passing it data, if any
	err = parsedTemplate.Execute(w, data)
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
		fmt.Fprint(w, "failed validation")
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	fmt.Println(email, password)

	fmt.Fprint(w, email)
}
