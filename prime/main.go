package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// n := -1
	// _, msg := isPrime(n)
	// fmt.Println(msg)

	// print welcome message
	intro()

	// create a channel to indicate when user wants to quit
	doneChannel := make(chan bool)

	// start a goroutine to read user input and run program
	go readUserInput(doneChannel)

	// block until doneChannel gets value
	<-doneChannel

	// close the channel
	close(doneChannel)

	// say goodbye
	fmt.Println("Goodbye.")
}

func readUserInput(doneChannel chan bool) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		res, done := checkNumber(scanner)

		if done {
			doneChannel <- true
			return
		}

		fmt.Println(res)
		prompt()
	}
}

func intro() {
	fmt.Println("Is it prime?")
	fmt.Println("------------")
	fmt.Println("Enter a whole number, and check if its a prime or not. Enter q to quit")
	prompt()
}

func checkNumber(scanner *bufio.Scanner) (string, bool) {
	// read uesr input
	scanner.Scan()

	// check if user wants to quit
	if strings.EqualFold(scanner.Text(), "q") {
		return "", true
	}

	// convert user input to int
	num, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return "Please enter a whole number", false
	}

	_, msg := isPrime(num)
	return msg, false
}

func prompt() {
	fmt.Print("-> ")
}

func isPrime(n int) (bool, string) {
	if n == 0 || n == 1 {
		return false, fmt.Sprintf("%d is not prime by definition", n)
	}

	if n < 0 {
		return false, "Negative numbers are not prime"
	}

	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false, fmt.Sprintf("%d is not prime because it is divisible by %d", n, i)
		}
	}

	return true, fmt.Sprintf("%d is a prime number", n)
}
