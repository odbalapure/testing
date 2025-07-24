package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFormHas(t *testing.T) {
	form := NewForm(nil)

	has := form.Has("age")
	if has {
		t.Error("age field is not part of the form")
	}

	postedData := url.Values{}
	postedData.Add("job", "Software Engineer")
	form = NewForm(postedData)

	has = form.Has("job")
	if !has {
		t.Errorf("field 'job' does not exist in the form")
	}
}

func TestFormRequired(t *testing.T) {
	r := httptest.NewRequest("POST", "/login", nil)
	form := NewForm(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r, _ = http.NewRequest("POST", "/login", nil)
	r.PostForm = postedData

	form = NewForm(r.PostForm)
	form.Required("a", "b", "c")

	if !form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}
}

func TestFormCheck(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")
	if form.Valid() {
		t.Error("Valid() returns false")
	}
}

func TestFormErrorGet(t *testing.T) {
	form := NewForm(nil)

	form.Check(false, "password", "password is required")
	s := form.Errors.Get("password")

	if len(s) == 0 {
		t.Error("should have an error returned from Get, but it does not")
	}

	s = form.Data.Get("job")
	if len(s) != 0 {
		t.Error("should not have an error, but got one")
	}
}
