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
		// Daemon work runs here:
		fmt.Println("Daemon is now running.")

		// Wait for kill signal to be received
		select {
		case <-sigChan:
			// Exit the program gracefully on signal
			fmt.Println("Kill signal received, stopping daemon.")
			os.Exit(0)
		}
	}
}
