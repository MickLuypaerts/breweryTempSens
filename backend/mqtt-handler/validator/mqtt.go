package validator

import (
	"fmt"
	"regexp"
	"strconv"
)

var ErrTooManyMatchesInTopic = fmt.Errorf("too many matches in topic")
var ErrTooFewMatchesInTopic = fmt.Errorf("not enough matches in topic")
var ErrGenericValidateTopic = fmt.Errorf("error validating topic")

func ValidateMQTTTopicAndRetID(topic string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	matches := re.FindAllString(topic, -1)
	if len(matches) == 1 {
		id, _ := strconv.Atoi(matches[0]) // should never fail because of regex
		return id, nil
	}
	if len(matches) > 1 {
		return -1, ErrTooManyMatchesInTopic
	}
	if len(matches) < 1 {
		return -1, ErrTooFewMatchesInTopic
	} else {
		return -1, ErrGenericValidateTopic
	}
}

func ValidateMQTTPayloadAndReturnFloat32(payload string) (float32, error) {
	temp, err := strconv.ParseFloat(payload, 32)
	if err != nil {
		return -1, err
	}
	return float32(temp), nil
}
