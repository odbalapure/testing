package main

import (
	"fmt"
	"net/http"
	"path"
	"text/template"
)

var pathToTempaltes = "./templates/"

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
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTempaltes, t))
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
