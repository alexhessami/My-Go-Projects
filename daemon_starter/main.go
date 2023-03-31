package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a channel to listen for signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Loop indefinitely
	for {
		// Do your daemon work here
		fmt.Println("Daemon is running...")

		// Wait for a signal to be received
		select {
		case <-sigChan:
			// Exit the program gracefully on signal
			fmt.Println("Received signal, stopping daemon.")
			os.Exit(0)
		}
	}
}
