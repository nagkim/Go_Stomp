package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-stomp/stomp/v3"
)

func main() {
	masterBroker := "10.37.129.2:61616"
	slaveBroker := "10.37.129.3:61616"
	queueName := "/test/message"
	username := "admin"
	password := "admin"

	connectToBroker := func(brokerURL string, username string, password string) (*stomp.Conn, error) {
		return stomp.Dial("tcp", brokerURL, stomp.ConnOpt.Login(username, password))
	}

	for {

		// Try connecting to the master broker first
		conn, err := connectToBroker(masterBroker, username, password)

		// If connection to the master fails, try connecting to the slave broker
		if err != nil {
			fmt.Println("Failed to connect to the master broker. Trying the slave broker...")
			conn, err = connectToBroker(slaveBroker, username, password)
		}

		if err != nil {
			fmt.Println("Failed to connect to ActiveMQ:", err)
			//	fmt.Println("Retrying in  seconds...")
			time.Sleep(time.Second) // Wait before attempting reconnection
			continue
		}

		defer conn.Disconnect()

		sub, err := conn.Subscribe(queueName, stomp.AckAuto)
		if err != nil {
			fmt.Println("Failed to subscribe to the queue:", err)
			return
		}
		defer sub.Unsubscribe()

		// Create a context and a cancel function for managing goroutines
		ctx, cancel := context.WithCancel(context.Background())

		// Receiving goroutine to handle incoming messages
		// Create an unbuffered channel to signal when receiving is done

		receiveMessages := make(chan struct{})
		go func() {
			// Infinite loop to keep listening for messages
			for {
				select {
				case msg := <-sub.C:
					// Check if the connection is closed
					if msg == nil {
						fmt.Println("Connection closed.")
						// Signal the main goroutine
						receiveMessages <- struct{}{}
						return
					}
					fmt.Println("Received message:", string(msg.Body))

				case <-ctx.Done():
					fmt.Println("Receiving goroutine canceled.")
					return
				}
			}
		}()

		// Waiting for the receiving goroutine to finish or connection loss
		<-receiveMessages
		// Cancel the receiving goroutine
		cancel()

	}
}
