package data

import (
	"mqtt-handler/validator"
)

type Reading struct {
	Temperture float32
	BarrelId   int
}

func NewReadingFromMQTT(payload, topic string) (*Reading, error) {
	temp, err := validator.ValidateMQTTPayloadAndReturnFloat32(payload)
	if err != nil {
		return nil, err
	}
	id, err := validator.ValidateMQTTTopicAndRetID(topic)
	if err != nil {
		return nil, err
	}
	return &Reading{temp, id}, nil
}

func NewReading(temp float32, id int) *Reading {
	return &Reading{temp, id}
}
