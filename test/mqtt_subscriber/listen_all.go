package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	// 1. Configuration
	broker := "tcp://148.230.101.142:1883"
	topic := "#"
	clientID := "fedora-laptop-client"

	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	// 2. Define the Message Handler
	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("📩 Received [%s]: %s\n", msg.Topic(), string(msg.Payload()))
	}

	opts.SetDefaultPublishHandler(messageHandler)

	// 3. Connection Callbacks
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("✅ Connected to MQTT broker!")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("❌ Connection lost: %v\n", err)
	}

	// 4. Create and Connect Client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("🔥 Error connecting: %v\n", token.Error())
		os.Exit(1)
	}

	// 5. Subscribe to Topic
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("🔥 Error subscribing: %v\n", token.Error())
		os.Exit(1)
	}

	fmt.Printf("🛰️  Listening on MQTT topic: %s\n", topic)

	// Keep the program running until you press Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n👋 Disconnecting...")
	client.Disconnect(250)
	time.Sleep(1 * time.Second)
}
