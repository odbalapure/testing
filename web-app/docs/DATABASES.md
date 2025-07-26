## Integrating databases

Install the postgres driver

> github.com/jackc/pgx/v4

Now write a function to connect to the psql database connection

```go
// db.go
import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	// Opened the connection but cannot ping
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	conn, err := openDB(app.DSN)

	if err != nil {
		return nil, err
	}

	log.Println("Connected to Postgres!")
	return conn, nil
}
```

Now open the database connection

```go
// main.go
type application struct {
	Session *scs.SessionManager
	DSN     string
	DB      *sql.DB
}

func main() {
    // ...
	app := application{}

    // when someone accesses `app.DSN` they get this though out their project
	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5432 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

    // When the main function exits, close the DB connection
	defer conn.Close()

	app.DB = conn
    // ...
}
```
