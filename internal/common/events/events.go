package events

import (
	"encoding/json"
	"reflect"
)

type EventMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EventData struct {
	EventId string
	Data    string
}

func ParseEventsResponse(inputValue string, res interface{}) (payload *EventData, err error) {
	payload = &EventData{}
	if err := json.Unmarshal([]byte(inputValue), payload); err != nil {
		return nil, err
	}

	result := reflect.ValueOf(res)
	err = json.Unmarshal([]byte(payload.Data), &result)
	return
}
