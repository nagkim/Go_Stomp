package main

import (
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

		fruits := []string{"Apple", "Banana", "Orange"}

		for _, v := range fruits {
			message := fmt.Sprintf("Message %s", v)
			err := conn.Send(queueName, "text/plain", []byte(message))
			if err != nil {
				fmt.Println("Failed to send message:", err)
				break
			}
			fmt.Println("Sent message:", message)
			time.Sleep(1 * time.Second) // Add a delay between sending messages
		}

		//fmt.Println("Retrying to connect to the servers...")
	}
}
