package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"mqtt-handler/data"
	"mqtt-handler/database"
	"mqtt-handler/mqtt_handler"
	"os"
	"os/signal"
	"syscall"
)

const (
	broker       = "broker.emqx.io"
	mqttPort     = 1883
	subTopic     = "brewery/fermentationbarrel/+/sensors/temperature"
	mqttUsername = "emqx"
	mqttPassword = "public"
	mqttClientId = "go_db_handler"

	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "brewery"
)

func main() {
	dbLogger := log.New(os.Stdout, "tempSens-db", log.LstdFlags)
	dbConn, err := database.NewDBConnection(dbHost, dbPort, dbUser, dbPassword, dbName, dbLogger)
	if err != nil {
		log.Fatal(err)
	}
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		reading, err := data.NewReadingFromMQTT(string(msg.Payload()), msg.Topic())
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Printf("id: %d\ttemp: %f\n", reading.BarrelId, reading.Temperture)
		dbConn.InsertProduct(reading)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	client, err := mqtt_handler.NewMQTTClient(broker, mqttPort, mqttClientId, mqttUsername, mqttPassword, messagePubHandler)
	if err != nil {
		log.Fatal(err)
	}

	client.Sub(subTopic)
	<-c
}
