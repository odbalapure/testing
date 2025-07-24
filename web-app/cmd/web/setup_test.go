package main

import (
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	app.Session = getSession()

	// m.Run() runs all the tests
	os.Exit(m.Run())
}
