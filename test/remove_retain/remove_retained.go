package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// Configure Client Options
	// Using the same broker as in listen_all.go
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://148.230.101.142:1883")
	opts.SetClientID("retained-remover")

	// Create and Connect the Client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("Connected to Broker")

	topicsToRemove := []string{
		"dbl/STM-1001/c",
		"dbl/SCM-1002/c",
	}

	for _, topic := range topicsToRemove {
		// Publish empty payload (nil) with retained=true to clear the message
		token := client.Publish(topic, 1, true, []byte{})
		token.Wait()
		if token.Error() != nil {
			fmt.Printf("Error clearing %s: %v\n", topic, token.Error())
		} else {
			fmt.Printf("Cleared retained message for: %s\n", topic)
		}
	}

	time.Sleep(500 * time.Millisecond)
	client.Disconnect(250)
	fmt.Println("Done.")
}
