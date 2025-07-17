package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestIsPrime(t *testing.T) {
	result, msg := isPrime(0)

	if result {
		t.Errorf("with %d as test parameter, go true, but expected false", 0)
	}

	if msg != "0 is not prime by definition" {
		t.Error("wrong message returned:", msg)
	}

	result, msg = isPrime(7)

	if !result {
		t.Errorf("with %d as test parameter, go true, but expected false", 0)
	}

	if msg != "7 is a prime number" {
		t.Error("wrong message returned:", msg)
	}
}

func TestIsPrimeV2(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number"},
		{"not prime", 8, false, "8 is not prime because it is divisible by 2"},
		{"zero", 0, false, "0 is not prime by definition"},
		{"one", 1, false, "1 is not prime by definition"},
		{"negative", -11, false, "Negative numbers are not prime"},
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

func TestPrompt(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	prompt()

	// close our writer
	_ = w.Close()

	// reset os.Stdout to what it was before
	os.Stdout = oldOut

	// read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	// perform our test
	if string(out) != "-> " {
		t.Errorf("incorrect prompt: expected -> but got %s", string(out))
	}
}

func TestIntro(t *testing.T) {
	// save a copy of os.Stdout
	oldOut := os.Stdout

	// create a read and write pipe
	r, w, _ := os.Pipe()

	// set os.Stdout to our write pipe
	os.Stdout = w

	intro()

	// close our writer
	_ = w.Close()

	// reset os.Stdout to what it was before
	os.Stdout = oldOut

	// read the output of our prompt() func from our read pipe
	out, _ := io.ReadAll(r)

	// perform our test
	if !strings.Contains(string(out), "Enter a whole number") {
		t.Errorf("intro text not correct: got %s", string(out))
	}
}
