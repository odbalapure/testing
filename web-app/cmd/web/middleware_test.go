package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddIpContext(t *testing.T) {
	tests := []struct {
		headerName  string
		headerValue string
		addr        string
		emptyAddr   bool
	}{
		{"", "", "", false},
		{"", "", "", true},
		{"X-Forwarded-For", "192.168.23.0", "", true},
		{"", "", "hello:world", true},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure the value exists in the context
		val := r.Context().Value(contextUserKey)
		if val == nil {
			t.Errorf("context value: %s not present", contextUserKey)
		}

		// make sure we got a string
		ip, ok := val.(string)
		if !ok {
			t.Error("the IP address context value is not a string")
		}
		t.Log(ip)
	})

	for _, e := range tests {
		// create a handler to test
		handlerToTest := app.addIPToContext(nextHandler)

		// create an HTTP request
		req := httptest.NewRequest("GET", "http://testing", nil)

		if e.emptyAddr {
			req.RemoteAddr = ""
		}

		if len(e.headerName) > 0 {
			req.Header.Add(e.headerName, e.headerValue)
		}

		if len(e.addr) > 0 {
			req.RemoteAddr = e.addr
		}

		handlerToTest.ServeHTTP(httptest.NewRecorder(), req)
	}
}

func TestAppIpFromContext(t *testing.T) {
	expectedIp := "0.0.0.0"

	ctx := context.WithValue(context.Background(), contextUserKey, expectedIp)
	ip := app.ipFromContext(ctx)

	if !strings.EqualFold(expectedIp, ip) {
		t.Errorf("expected %s got %s", expectedIp, ip)
	}
}
