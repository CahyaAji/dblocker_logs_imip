package service

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	Client mqtt.Client
}

func NewMqttClient(broker string, clientID string) (*MqttClient, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetAutoReconnect(true)

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("Connected to MQTT broker")
	})
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Printf("Connection lost: %v", err)
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		log.Println("Attempting to reconnect to MQTT broker...")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &MqttClient{Client: client}, nil
}

func (m *MqttClient) Publish(topic string, payload interface{}) error {
	token := m.Client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func (m *MqttClient) Unsubscribe(topics ...string) error {
	token := m.Client.Unsubscribe(topics...)
	token.Wait()
	return token.Error()
}

func (m *MqttClient) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := m.Client.Subscribe(topic, 0, handler)
	token.Wait()
	return token.Error()
}
