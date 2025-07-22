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