#include "barrel.h"
#include "DHT.h"

Barrel createNewBarrel(int id, int pin, uint8_t dhtType, const char* mqttPubTopicFormat) {
	Barrel barrel;
	barrel.ID = id;
	barrel.enabled = true;
	barrel.temp = 0;	
	barrel.millis = 0;
	barrel.dht = new DHT(pin, dhtType);
	barrel.dht->begin();
	sprintf(barrel.mqttPubTopic, mqttPubTopicFormat, barrel.ID);
	return barrel;
}