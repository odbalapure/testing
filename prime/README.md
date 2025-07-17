## Testing

Create a go app

```go
go mod init <project-name>
```

Now create a test file with the name `file_name_test.go`.

> Go has a convention of placing the test file and source code side by side.

And the test functions start with `Test`, eg: `TestPrime`. The test function name can be `TestIsPrime` or `Test_isPrime`. Using something like `TestisPrime` will result in an error:

```go
./main_test.go:7:6: TestisPrime has malformed name: first letter after 'Test' must not be lowercase
FAIL    prime [build failed]
FAIL
```

## Writing a simple test

```go
func TestIsPrime(t *testing.T) {
	result, msg := isPrime(0)

	if result {
		t.Errorf("with %d as test parameter, go true, but expected false", 0)
	}

	if msg != "0 is not prime!" {
		t.Error("wrong message returned:", msg)
	}

    	result, msg = isPrime(7)

	if !result {
		t.Errorf("with %d as test parameter, go true, but expected false", 0)
	}

	if msg != "7 is a prime number!" {
		t.Error("wrong message returned:", msg)
	}
}
```

**NOTE**: As you can we doing lot of repition; we can use table tests to deal with this.


## Table tests

This will let us test multiple edge cases and conditions without duplication.

```go
func TestIsPrimeV2(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number"},
		{"not prime", 8, false, "8 is not prime because it is divisible by 2"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)

		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}
```

Coverage can be checked using

```go
go test -cover .
```

Also coverage details can be written to a file; check the coverage result in the browser. This will open the coverage report in the browser

```go
go test -coverprofile=coverag.out
go tool cover -html=coverage.out
```
 
Where `-coverprofile=coverag.out` saves the coverage result to a file called `coverage.out`.
