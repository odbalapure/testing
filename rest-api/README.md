## REST API

- [DB functions](./cmd/api/db.go)
- [Setting up routes](./cmd/api/routes.go)
- [Starting the app and db server](./cmd/api/main.go)

### Creating a REST endpoint in Go

Inside the routes.go file, add a new route the following way

```go
mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
    var payload = struct {
        Message string `json:"message"`
    }{
        Message: "Hello world!",
    }

    _ = app.writeJSON(w, http.StatusOK, payload)
})
```

The writeJSON is a parses the data passed to it into JSON format


```go
var out []byte
// ...
jsonBytes, err := json.Marshal(data)
if err != nil {
    return err
}
out = jsonBytes

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)

_, err := w.Write(out)
if err != nil {
    return err
}
return nil
// ...
```
