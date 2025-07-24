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

### Installing a session package
