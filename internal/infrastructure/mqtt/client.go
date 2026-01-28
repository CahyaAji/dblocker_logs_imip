package mqtt

import (
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Client interface {
	Publish(topic string, qos byte, retained bool, payload any) error
	Subscribe(topic string, qos byte, handler paho.MessageHandler) error
	Unsubscribe(topics ...string) error
	Close()
}

type mqttClient struct {
	pahoClient paho.Client
}

func New(broker string, clientID string) (Client, error) {
	opts := paho.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetAutoReconnect(true)

	opts.OnConnect = func(c paho.Client) {
		log.Println("Connected to MQTT broker")

	}
	opts.OnConnectionLost = func(c paho.Client, err error) {
		log.Printf("Connection lost: %v", err)
	}

	opts.OnReconnecting = func(c paho.Client, options *paho.ClientOptions) {
		log.Println("Attempting to reconnect to MQTT broker...")
	}

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt connect error: %w", token.Error())
	}

	return &mqttClient{pahoClient: client}, nil

}

func (m *mqttClient) Publish(topic string, qos byte, retained bool, payload any) error {
	token := m.pahoClient.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}

func (m *mqttClient) Subscribe(topic string, qos byte, handler paho.MessageHandler) error {
	token := m.pahoClient.Subscribe(topic, qos, handler)
	token.Wait()
	return token.Error()
}

func (m *mqttClient) Unsubscribe(topics ...string) error {
	token := m.pahoClient.Unsubscribe(topics...)
	token.Wait()
	return token.Error()
}

func (m *mqttClient) Close() {
	m.pahoClient.Disconnect(250)
}
