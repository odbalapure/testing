package main

import (
	"fmt"
	"net/http"
	"text/template"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "Home Page")
	err := app.render(w, r, "home.page.gohtml", &TemplateData{})
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
	parsedTemplate, err := template.ParseFiles("./templates/" + t)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	// execute the template, passing it data, if any
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}
