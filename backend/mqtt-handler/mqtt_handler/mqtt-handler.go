package mqtt_handler

import (
	"fmt"
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

type MQTTHandler struct {
	Client mqtt.Client
}

func NewMQTTClient(broker string, port int, clientId string, username string, userPass string, messagePubHandler func(client mqtt.Client, msg mqtt.Message)) (MQTTHandler, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientId)
	opts.SetUsername(username)
	opts.SetPassword(userPass)
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return MQTTHandler{}, token.Error()
	}
	return MQTTHandler{client}, nil
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v\n", err)
}

func (c MQTTHandler) Sub(topic string) {
	token := c.Client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic: %s\n", topic)
}
