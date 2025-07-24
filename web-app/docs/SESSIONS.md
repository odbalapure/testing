## Testing sessions

Sessions keep state for web apps that render pages on the server side.

Create a setup_test.go file.

> Make sure the name remains the same

```go
import (
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	// m.Run() runs all the tests
	os.Exit(m.Run())
}
```

`testing.M` runs before any of the tests in our code base. And `m.Run()` runs all the tests.

We need this setup because we end up duplicating test code. Eg: We can now remove `var app application` from all the test suites.

### Creating a sesssion

Install a session package

```sh
go get github.com/alexedwards/scs/v2
```

Creating a session

```go
// session.go
import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

func getSession() *scs.SessionManager {
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

// handler.go
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var td = make(map[string]any)

	if app.Session.Exists(r.Context(), "test") {
		msg := app.Session.Get(r.Context(), "test")
		td["test"] = msg
	} else {
		app.Session.Put(r.Context(), "test", "Hit this page "+time.Now().UTC().String())
	}

	err := app.render(w, r, "home.page.gohtml", &TemplateData{Data: td})
	fmt.Print(err)
}

// routes.go
mux.Use(app.Session.LoadAndSave)

// main.go
func main() {
	app := application{}

	// get a session manager
	// NOTE: setup the session before hitting the routes
	app.Session = getSession()

    mux := app.routes()

	log.Println("Strating serve on port 8080")

	// start the server
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// setup_test.go
func TestMain(m *testing.M) {
    // Make sure to create a session and make it accessible to all the tests
    // Otherwise it result in the following error:
    // panic: runtime error: invalid memory address or nil pointer dereference
	app.Session = getSession()

	os.Exit(m.Run())
}
```

## Testing controller

```go
func TestAppHom(t *testing.T) {
	// create a request
	req, _ := http.NewRequest("GET", "/", nil)

	req = addContextAndSessionToRequest(req, app)

	// `NewRecorder()` will return the following:
	// {
	// 	HeaderMap: make(http.Header),
	// 	Body:      new(bytes.Buffer),
	// 	Code:      200,
	// }
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(app.Home)
	handler.ServeHTTP(rr, req)

	// check status code
	if rr.Code != http.StatusOK {
		t.Errorf("TestAppHome returned wrong status code; expected 200 got %d", rr.Code)
	}
}
```

> Running the `main` function test as it is; will result in `panic: scs: no session data in context`. So, we need to setup a "contextUserkey" and the "session".

```go
func getCtx(req *http.Request) context.Context {
	ctx := context.WithValue(req.Context(), contextUserKey, "unknown")
	return ctx
}

func addContextAndSessionToRequest(req *http.Request, app application) *http.Request {
	req = req.WithContext(getCtx(req))

    // the test won't have an X-Session header unless we explicitly add it
	ctx, _ := app.Session.Load(req.Context(), req.Header.Get("X-Session"))

	return req.WithContext(ctx)
}
```
