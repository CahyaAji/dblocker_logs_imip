package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func decodeAndPrintCommand(topic string, payload []byte) {
	if len(payload) != 2 {
		fmt.Printf("CMD payload invalid length: %d\n", len(payload))
		fmt.Printf("Payload: %x\n", payload)
		return
	}

	mask := uint16(payload[0])<<8 | uint16(payload[1])

	fmt.Println("───── DBlocker Command ─────")
	fmt.Printf("Topic    : %s\n", topic)
	fmt.Printf("Raw mask : 0x%04X\n", mask)
	fmt.Printf("Binary   : %016b\n", mask)
	fmt.Println("States:")

	print := func(name string, bit int) {
		state := "OFF"
		if mask&(1<<bit) != 0 {
			state = "ON"
		}
		fmt.Printf("  %-16s : %s\n", name, state)
	}

	print("SignalGPS[0]", 0)
	print("SignalCtrl[0]", 1)

	print("SignalGPS[1]", 2)
	print("SignalCtrl[1]", 3)

	print("SignalGPS[2]", 4)
	print("SignalCtrl[2]", 5)

	print("FanMaster", 6)

	print("SignalGPS[3]", 7)
	print("SignalCtrl[3]", 8)

	print("SignalGPS[4]", 9)
	print("SignalCtrl[4]", 10)

	print("SignalGPS[5]", 11)
	print("SignalCtrl[5]", 12)

	print("FanSlave", 13)

	fmt.Println("───────────────────────────")
}

// 1. Define the Message Handler
// This function is called whenever a message is received.
var messagePubHandler mqtt.MessageHandler = func(
	client mqtt.Client,
	msg mqtt.Message,
) {
	topic := msg.Topic()
	parts := strings.Split(topic, "/")

	// Expect: dbl/{serial}/c
	if len(parts) >= 3 && parts[2] == "c" {
		decodeAndPrintCommand(topic, msg.Payload())
		return
	}

	// Other topics
	fmt.Printf("Topic: %s\n", topic)
	fmt.Printf("Payload: %s\n", string(msg.Payload()))
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
