package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		fmt.Printf("Error on creating connection: %s", err)
		return
	}
	fmt.Println("Connection established")
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("Error on creating connection channel: %s", err)
		return
	}

	defer func() {
		fmt.Println("Closing AMQP connection...")
		if err := conn.Close(); err != nil {
			fmt.Printf("Error closing connection: %s", err)
		}
	}()

	arg := routing.PlayingState{}
	arg.IsPaused = true
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, arg)
	if err != nil {
		fmt.Printf("Error publishing json: %s", err)
	}

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	<-c
	fmt.Println("Got signal")
}
