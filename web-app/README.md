## Testing web app

### Creating a simple web server

This is the `main.go` code

```go
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
```

### Creating routes

Create a `routes.go` file; add the routes and middleware code.

```go
import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// middleware
	mux.Use(middleware.Recoverer)

	// register routes
	mux.Get("/v1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// static assets

	return mux
}
```

> Start the server using `go run ./cmd/web` because the templates folder is present in the "root". Executing it from nested folder will result in an error; because the CWD is not "root". 

### Writing tests for our routes

```go
package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAppRoutes(t *testing.T) {
	var registered = []struct {
		route  string
		method string
	}{
		{"/", "GET"},
	}

	var app application
	// Returns type `http.Handler`
	mux := app.routes()

	// Casting chi.Routes
	chiRoutes := mux.(chi.Routes)

	for _, route := range registered {
		// check if the route exists
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("route %s is not registered", route.route)
		}
	}
}

func routeExists(testRoute string, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			found = true
		}

		return nil
	})

	return found
}
```

### Writing a middleware

Middleware is a function that sits b/w incoming request and main application logic. Eg: A middleware that extracts an IP and adds it to the context before it reaches a route.

```go
import (
	"context"
	"fmt"
	"net"
	"net/http"
)

type contextKey string

const contextUserKey contextKey = "user_ip"

func (app *application) ipFromContext(ctx context.Context) string {
	return ctx.Value(contextUserKey).(string)
}

func (app *application) addIPToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx = context.Background()
		// get the ip
		ip, err := getIP(r)
		if err != nil {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
			if len(ip) == 0 {
				ip = "unknown"
			}
			ctx = context.WithValue(r.Context(), contextUserKey, ip)
		} else {
			ctx = context.WithValue(r.Context(), contextUserKey, ip)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "unknown", err
	}

	// check if an IP is valid
	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "", fmt.Errorf("the user ip:%s is not in the format of ip:port", r.RemoteAddr)
	}

	// is it coming from a proxy
	forward := r.Header.Get("X-Forwarded-For")
	if len(forward) > 0 {
		ip = forward
	}

	if len(ip) == 0 {
		ip = "forward"
	}

	return ip, nil
}
```

Accessing the IP address from the context 

```go
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
```
