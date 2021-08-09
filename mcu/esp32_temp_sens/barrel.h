#ifndef BARREL_H
#define BARREL_H

#include <stdbool.h>
#include "DHT.h"

typedef struct {
    int ID;
    int pin;
    float temp;
    bool enabled;
    DHT* dht;
    unsigned long millis;
    char mqttPubTopic[30];
} Barrel;

Barrel createNewBarrel(int id, int pin, uint8_t dhtType, const char* mqttPubTopicFormat);

#endif