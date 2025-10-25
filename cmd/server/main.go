package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close()
	fmt.Println("Connected to RabbitMQ successfully")

	// Create a channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	// Notify the channel on SIGINT (Ctrl+C) or SIGTERM
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create a channel to signal program exit
	done := make(chan bool, 1)

	// Goroutine to handle the signal - This happend concurrently with the code after the function
	go func() {
		<-signalChan // Wait for a signal
		fmt.Println("\nReceived an interrupt, stopping gracefully...")
		done <- true
	}()

	fmt.Println("Program is running. Press Ctrl+C to stop.")
	<-done // Block until done signal is received
	fmt.Println("Program stopped.")

}
