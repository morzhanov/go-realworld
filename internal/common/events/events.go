package events

import (
	"encoding/json"
	"reflect"
)

type EventData struct {
	EventId string
	Data    string
}

func ParseEventsResponse(inputValue []byte, res interface{}) (payload *EventData, err error) {
	payload = &EventData{}
	if err := json.Unmarshal([]byte(inputValue), payload); err != nil {
		return nil, err
	}

	result := reflect.ValueOf(res)
	err = json.Unmarshal([]byte(payload.Data), &result)
	return
}
