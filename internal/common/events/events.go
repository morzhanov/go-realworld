package events

import (
	"encoding/json"
)

type EventData struct {
	EventId string
	Data    string
}

func ParseEventsResponse(inputValue []byte, res interface{}) (payload *EventData, err error) {
	payload = &EventData{}
	if err := json.Unmarshal(inputValue, payload); err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(payload.Data), &res)
	return
}
