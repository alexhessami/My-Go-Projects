package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getNumber() int {
	randomValue := rand.Intn(100)
	return randomValue

}

func ask() int {
	var input int
	fmt.Print("I am thinking of a number between 1-100: ")
	// Get the input from the user
	fmt.Scan(&input)
	return input
}

func main() {

	rand.Seed(time.Now().UnixNano())

	var guess int
	guess = ask()
	number := getNumber()
	tries := 0

	for {

		if tries == 6 {
			fmt.Println("Sorry you have ran out of guesses.")
			fmt.Println("The number was", number, ".")
			break
		}

		if guess > number {
			tries += 1
			fmt.Println("You have ", 6-tries, "tries left.")
			if tries < 6 {
				fmt.Println("The number is smaller.")
				guess = ask()
			}
		} else if guess < number {
			tries += 1
			fmt.Println("You have ", 6-tries, "tries left.")
			if tries < 6 {
				fmt.Println("The number is bigger.")
				guess = ask()
			}
		} else if guess == number {
			tries += 1
			fmt.Println("You are correct! You may go through to the treasure!")
			fmt.Println("You won in ", tries, "tries!")
			break
		}

	}
}
