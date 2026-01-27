package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// 1. Define the Message Handler
// This function is called whenever a message is received.
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected to Broker")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

func main() {

	// 2. Configure Client Options
	// Using a public test broker for demonstration.
	// Change "tcp://broker.hivemq.com:1883" to your specific broker URL if needed.
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://148.230.101.142:1883")
	opts.SetClientID("listener-all")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// 3. Create and Connect the Client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// 4. Subscribe to All Topics
	// The "#" wildcard matches any topic.
	// QoS 1 means "at least once" delivery.
	topic := "#"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)

	// 5. Keep the Program Running
	// This block waits for a CTRL+C (SIGINT) signal to stop the program gracefully.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	fmt.Println("Disconnecting...")
	client.Disconnect(250)
}
