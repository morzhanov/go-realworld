package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/morzhanov/go-realworld/internal/users/dto"
	"github.com/morzhanov/go-realworld/internal/users/services"
	"github.com/segmentio/kafka-go"
)

type UsersEventsController struct {
	service *services.UsersService
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

type GetUserDataInput struct {
	BaseEventPayload
	UserId string `json:"userId"`
}

type GetUserDataByUsernameInput struct {
	BaseEventPayload
	Username string `json:"username"`
}

type ValidateUserPasswordInput struct {
	BaseEventPayload
	dto.ValidateUserPasswordDto
}

type CreateUserInput struct {
	BaseEventPayload
	dto.CreateUserDto
}

type DeleteUserInput struct {
	BaseEventPayload
	UserId string `json:"userId"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}

type SuccessMessage struct{}

func check(err error) {
	// TODO: handle error
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: topic name should be service name
func CreateKafkaConnection(topic string, partition int) *kafka.Conn {
	// TODO: provide kafka uri
	uri := "192.168.0.180:32181"
	conn, _ := kafka.DialLeader(context.Background(), "tcp", uri, topic, partition)
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	return conn
}

func (c *UsersEventsController) Listen() {
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
func (c *UsersEventsController) processRequest(b *[]byte) {
	input := EventMessage{}
	err := json.Unmarshal(*b, &input)
	check(err)

	switch input.Key {
	case "users:get_data":
		payload := GetUserDataInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.GetUserData(payload.UserId)
		check(err)

		c.sendResponse(payload.EventId, res)
	case "users:get_data_by_username":
		payload := GetUserDataByUsernameInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.GetUserDataByUsername(payload.Username)
		check(err)

		c.sendResponse(payload.EventId, res)
	case "users:validate_password":
		payload := ValidateUserPasswordInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		err = c.service.ValidateUserPassword(&payload.ValidateUserPasswordDto)
		if err != nil {
			c.sendResponse(payload.EventId, &ErrorMessage{Error: err.Error()})
		}
		c.sendResponse(payload.EventId, &SuccessMessage{})
	case "users:create":
		payload := CreateUserInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		res, err := c.service.CreateUser(&payload.CreateUserDto)
		check(err)

		c.sendResponse(payload.EventId, res)
	case "users:delete":
		payload := DeleteUserInput{}
		err := json.Unmarshal([]byte(input.Value), &payload)
		check(err)

		if err != nil {
			c.sendResponse(payload.EventId, &ErrorMessage{Error: err.Error()})
		}
		c.sendResponse(payload.EventId, &SuccessMessage{})
	}

}

// TODO: send the response with client module and use payload.EventId as event_uuid
func (c *UsersEventsController) sendResponse(eventUuid string, value interface{}) {
	// TODO: send response via client kafka
	payload, err := json.Marshal(&value)
	check(err)
	fmt.Printf("payload %v\n", payload)
}

func NewUsersEventsController(s *services.UsersService) *UsersEventsController {
	// TODO: provide topic and partition from config
	conn := CreateKafkaConnection("topic", 0)
	return &UsersEventsController{service: s, Conn: conn}
}
