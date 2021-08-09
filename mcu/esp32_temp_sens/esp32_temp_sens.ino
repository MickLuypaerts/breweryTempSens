#include "secrets.h"
#include "barrel.h"
#include "DHT.h"
#include <WiFi.h>
#include <PubSubClient.h>

#define BAUD_RATE 115200

#define DHTTYPE DHT11
#define DHT_PIN_0 15
#define DHT_PERIOD 2000

#define MQTT_SERVER "test.mosquitto.org"
#define MQTT_SUB_TOPIC "brewery/fermentationbarrel/+/controls/temperaturesensor"
#define MQTT_PUB_TOPIC_FORMAT "brewery/fermentationbarrel/%d/sensors/temperature"
#define MQTT_PORT 1883
#define MQTT_RECONNECT_DELAY 5000

void getDhtData(Barrel* barrel);
void connectWifi();
void checkConnection();
void reconnectMQTT();
void mqttCallback(char* topic, byte* payload, unsigned int length);

#define BARREL_LENGHT 1
Barrel barrels[BARREL_LENGHT] = {createNewBarrel(0, DHT_PIN_0, DHTTYPE, MQTT_PUB_TOPIC_FORMAT)};

WiFiClient espClient;
PubSubClient mqttClient(espClient);


void setup() {
    Serial.begin(BAUD_RATE);
    Serial.println("Starting");
    connectWifi();
    mqttClient.setServer(MQTT_SERVER, MQTT_PORT);
    mqttClient.setCallback(mqttCallback);
    reconnectMQTT();
}

void loop() {
    // for each barrel see if it's enabled if so get the temp (if more than 2 second age) en send this to the mqtt broker
    for(int i = 0; i < BARREL_LENGHT; i++) {
        if (barrels[i].enabled) {
            if (millis() - barrels[i].millis >= DHT_PERIOD) {
            getDhtData(&barrels[i]);
            char temp[5];
            sprintf(temp, "%.2f", barrels[i].temp);
            mqttClient.publish(barrels[i].mqttPubTopic, temp);
            barrels[i].millis = millis();
            }
        }
    }
    checkConnection();
    mqttClient.loop();
}

void getDhtData(Barrel* barrel) {
  barrel->temp = barrel->dht->readTemperature();
  if (isnan(barrel->temp)) {
    Serial.printf("getDhtData: Failed to read from DHT sensor %d\n", barrel->ID);
    return;
  }
  Serial.println(barrel->temp);
}

void connectWifi() {
  Serial.print("Connecting to Wifi .");
  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASS);
  while (WiFi.status() != WL_CONNECTED) {
    Serial.print(".");
    delay(500);
  }
  Serial.println("\nConnected to Wifi.");
}

void checkConnection() {
  if (WiFi.status() != WL_CONNECTED) {
    Serial.println("WiFi disconnected reconnecting.");
    connectWifi();
    reconnectMQTT();
  } if (!mqttClient.connected()) {
    reconnectMQTT();
  }
}

void reconnectMQTT() {
  // Loop until we're reconnected
  while (!mqttClient.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Attempt to connect
    if (mqttClient.connect("esp32Client_ferm_temp")) {
      Serial.println("connected");
      // ... and resubscribe
      mqttClient.subscribe(MQTT_SUB_TOPIC);
    } else {
      Serial.print("failed, rc=");
      Serial.print(mqttClient.state());
      Serial.println(" try again in 5 seconds");
      delay(MQTT_RECONNECT_DELAY);
    }
  }
}

void mqttCallback(char* topic, byte* payload, unsigned int length) {
    if(length != 1) {
      return;
    }
    int value = (char)payload[0] - '0';
    int id = (char)topic[27] - '0';

    for(int i = 0; i < BARREL_LENGHT; i++) {
		if(barrels[i].ID == id ) {
			if(value == 1) {
				barrels[i].enabled = true;
				Serial.printf("Barrel %d enabled\n", barrels[i].ID);
				return;
			} else if(value == 0) {
				barrels[i].enabled = false;
				Serial.printf("Barrels %d disabled\n", barrels[i].ID);
			}
		}
    }
}
