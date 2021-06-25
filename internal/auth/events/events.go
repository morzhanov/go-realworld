package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/morzhanov/go-realworld/internal/auth/dto"
	"github.com/morzhanov/go-realworld/internal/auth/services"
	"github.com/segmentio/kafka-go"
)

type AuthEventsController struct {
	service *services.AuthService
	Conn    *kafka.Conn
}

// TODO: move to common package
type EventMessage struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type BaseEventPayload struct {
	EventId string `json:"eventId"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct{}

type LoginInput struct {
	BaseEventPayload
	dto.LoginInput
}

type Signup struct {
	BaseEventPayload
	dto.SignupInput
}

type ValidateEventsRequestInput struct {
	BaseEventPayload
	dto.ValidateEventsRequestInput
}

func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: looks like common function for all controllers
// TODO: topic name should be service name
func CreateKafkaConnection(topic string, partition int) *kafka.Conn {
	// TODO: provide kafka uri
	uri := "192.168.0.180:32181"
	conn, _ := kafka.DialLeader(context.Background(), "tcp", uri, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return conn
}

// TODO: looks like common function for all controllers
func (c *AuthEventsController) Listen() {
	c.Conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	batch := c.Conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		_, err := batch.Read(b)
		if err != nil {
			break
		}
		go c.processRequest(&b)
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}
	if err := c.Conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}
}

/*
Request event schema:
	Key: "resource:action"
	Value: "{payload}"
*/
/*
Response event schema:
	Key: "result:event_uuid"
	Value: "{payload}"
*/
func (c *AuthEventsController) processRequest(b *[]byte) {
	input := EventMessage{}
	err := json.Unmarshal(*b, &input)
	check(err)

	switch input.Key {
	case "auth:validate_events_request":
		payload := ValidateEventsRequestInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		err = c.service.ValidateEventsRequest(&payload.ValidateEventsRequestInput)
		if err != nil {
			c.sendResponse(payload.EventId, &ErrorMessage{Error: err.Error()})
		}
		c.sendResponse(payload.EventId, &SuccessMessage{})
	case "auth:login":
		payload := LoginInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.Login(&payload.LoginInput)
		check(err)

		c.sendResponse(payload.EventId, res)
	case "pictures:signup":
		payload := Signup{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.Signup(&payload.SignupInput)
		check(err)
		c.sendResponse(payload.EventId, res)
	}
}

// TODO: send the response with client module and use payload.EventId as event_uuid
func (c *AuthEventsController) sendResponse(eventUuid string, value interface{}) {
	// TODO: send response via client kafka
	payload, err := json.Marshal(&value)
	check(err)
	fmt.Printf("payload %v\n", payload)
}

func NewAuthEventsController(s *services.AuthService) *AuthEventsController {
	// TODO: provide topic and partition from config
	conn := CreateKafkaConnection("topic", 0)
	return &AuthEventsController{service: s, Conn: conn}
}
